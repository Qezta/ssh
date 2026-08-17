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
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	gossh "github.com/charmbracelet/ssh"

	"github.com/Qezta/ssh/internal/tui"
)

func main() {
	host := envOr("SSH_HOST", "0.0.0.0")
	port := envOr("SSH_PORT", "23234")
	hostKeyPath := envOr("SSH_HOST_KEY_PATH", ".ssh/qezta_ed25519")

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(hostKeyPath),
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

	m := tui.NewModel(s.User(), pty.Window.Width, pty.Window.Height)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
