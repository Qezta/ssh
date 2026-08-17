// Package tui implements the terminal UI shown to every SSH connection.
// It mirrors the fake browser terminal on divit.qezta.com
// (src/lib/components/Terminal.svelte + src/lib/data/terminal-commands.ts),
// but backed by a real SSH session instead of a JS widget.
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logo = `      _ _       _ _
  ___(_) |_ ___| | |
 / __| | __/ _ \ | |
| (__| | ||  __/_|_|
 \___|_|\__\___(_|_)`

const bannerDivit = logo + `

divit@qezta.com — type 'help' to get started, Ctrl+C to disconnect.`

const bannerQezta = logo + `

qezta.com — under active development. Type 'help' to get started,
Ctrl+C to disconnect. Looking for Divit? ssh divit@qezta.com`

var (
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	bannerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("117"))
)

// Model is the bubbletea model driving one SSH session's terminal.
type Model struct {
	username string
	identity identity
	width    int
	height   int

	input   textinput.Model
	lines   []string
	history []string
	histPos int
}

// NewModel builds the initial state for a freshly connected session.
// username is whatever the client passed to `ssh <username>@qezta.com`; it
// picks which persona (identity) the session sees — see identityFor.
func NewModel(username string, width, height int) Model {
	if username == "" {
		username = "guest"
	}
	id := identityFor(username)

	ti := textinput.New()
	ti.Prompt = prompt(username)
	ti.Focus()
	ti.CharLimit = 256

	banner := bannerQezta
	if id == identityDivit {
		banner = bannerDivit
	}

	return Model{
		username: username,
		identity: id,
		width:    width,
		height:   height,
		input:    ti,
		lines:    []string{bannerStyle.Render(banner), ""},
	}
}

func prompt(username string) string {
	return promptStyle.Render(fmt.Sprintf("%s@qezta:~$ ", username))
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.input.Width = m.width - lipgloss.Width(m.input.Prompt) - 1
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			return m.runLine(), nil

		case tea.KeyUp:
			m.recallHistory(-1)
			return m, nil

		case tea.KeyDown:
			m.recallHistory(1)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Model) recallHistory(delta int) {
	if len(m.history) == 0 {
		return
	}
	m.histPos += delta
	if m.histPos < 0 {
		m.histPos = 0
	}
	if m.histPos >= len(m.history) {
		m.histPos = len(m.history)
		m.input.SetValue("")
		return
	}
	m.input.SetValue(m.history[m.histPos])
	m.input.CursorEnd()
}

func (m Model) runLine() Model {
	line := m.input.Value()
	m.input.SetValue("")

	trimmed := strings.TrimSpace(line)
	m.lines = append(m.lines, m.input.Prompt+line)

	if trimmed != "" {
		m.history = append(m.history, trimmed)
	}
	m.histPos = len(m.history)

	result := executeCommand(m.identity, line)
	switch {
	case result.clear:
		m.lines = nil
	case result.isError:
		m.lines = append(m.lines, errorStyle.Render(result.output), "")
	case result.output != "":
		m.lines = append(m.lines, result.output, "")
	default:
		m.lines = append(m.lines, "")
	}

	return m
}

func (m Model) View() string {
	var b strings.Builder

	visible := m.lines
	// Keep the scrollback from growing the rendered frame past the
	// session's terminal height (leave room for the input line).
	maxLines := m.height - 2
	if maxLines > 0 && len(visible) > maxLines {
		visible = visible[len(visible)-maxLines:]
	}

	for _, line := range visible {
		b.WriteString(line)
		b.WriteString("\n")
	}
	b.WriteString(m.input.View())

	return b.String()
}
