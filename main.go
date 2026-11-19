// Command ssh runs the SSH server behind `ssh divit@qezta.com`.
//
// It never grants a real shell: every connection — any username, any key or
// password — is dropped straight into a sandboxed Bubble Tea TUI (see
// internal/tui) that mirrors the fake terminal on divit.qezta.com.
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	gossh "github.com/charmbracelet/ssh"
	"github.com/muesli/termenv"

	"github.com/Qezta/ssh/internal/tui"
)

func main() {
	host := envOr("SSH_HOST", "0.0.0.0")
	port := envOr("SSH_PORT", "23234")
	hostKeyPath := envOr("SSH_HOST_KEY_PATH", ".ssh/qezta_ed25519")

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(hostKeyPath),
		// Anyone on the internet can connect, so cap how long a session may
		// hold resources: idle tabs and abandoned connections must not
		// accumulate until the pod hits its memory limit.
		wish.WithIdleTimeout(15*time.Minute),
		wish.WithMaxTimeout(2*time.Hour),
		// This is a public novelty terminal, not real shell access: accept
		// every username/key/password combination rather than gating auth.
		wish.WithPublicKeyAuth(func(ctx gossh.Context, key gossh.PublicKey) bool {
			return true
		}),
		wish.WithPasswordAuth(func(ctx gossh.Context, password string) bool {
			return true
		}),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	log.Printf("starting ssh server on %s:%s", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, gossh.ErrServerClosed) {
			log.Fatalf("could not start server: %v", err)
		}
	}()

	<-done
	log.Println("stopping ssh server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func teaHandler(s gossh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}

	m := tui.NewModel(s.User(), sessionRenderer(s), pty.Window.Width, pty.Window.Height)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

// sessionRenderer picks the client's color profile from its TERM and
// environment alone. wish's MakeRenderer also queries the terminal's
// background color; when a client doesn't answer, the reader that query leaves
// behind swallows the session's first keystrokes, and every connection waits
// out its timeout. The TUI uses no adaptive colors, so it never needs the
// answer.
func sessionRenderer(s gossh.Session) *lipgloss.Renderer {
	pty, _, _ := s.Pty()
	env := sessionEnv(append(s.Environ(), "TERM="+pty.Term))
	// WithUnsafe: the session is not a TTY on the server side, but the client
	// asked for a PTY, so its environment describes a real terminal.
	return lipgloss.NewRenderer(s, termenv.WithEnvironment(env), termenv.WithUnsafe(), termenv.WithColorCache(true))
}

// sessionEnv exposes the client's environment to termenv.
type sessionEnv []string

func (e sessionEnv) Environ() []string { return e }

// Getenv returns the last value set for key, so the PTY's TERM, appended
// after the client's own variables, wins.
func (e sessionEnv) Getenv(key string) string {
	for i := len(e) - 1; i >= 0; i-- {
		if v, ok := strings.CutPrefix(e[i], key+"="); ok {
			return v
		}
	}
	return ""
}
