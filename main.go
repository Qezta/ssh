// Command ssh runs the SSH server behind `ssh divit@qezta.com`.
//
// It never grants a real shell. divit@/div@ and hi@/hello@ get a sandboxed
// Bubble Tea TUI (see internal/tui); any other username gets a short guide
// to those entry points and is disconnected.
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

	tea "charm.land/bubbletea/v2"
	gossh "charm.land/ssh"
	"charm.land/wish/v2"
	bm "charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/colorprofile"

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
	if guide := tui.Guide(s.User()); guide != "" {
		wish.Println(s, strings.ReplaceAll(guide, "\n", "\r\n"))
		_ = s.Exit(0)
		return nil, nil
	}

	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}

	m := tui.NewModel(s.User(), pty.Window.Width, pty.Window.Height)
	return m, []tea.ProgramOption{tea.WithColorProfile(sessionProfile(s))}
}

// sessionProfile picks the client's color profile from its TERM and
// environment alone, without querying the client terminal. Some clients
// never answer terminal queries, which used to stall the first keystrokes.
// The PTY's TERM is appended last so it wins over the client's own variables.
func sessionProfile(s gossh.Session) colorprofile.Profile {
	pty, _, _ := s.Pty()
	return colorprofile.Env(append(s.Environ(), "TERM="+pty.Term))
}
