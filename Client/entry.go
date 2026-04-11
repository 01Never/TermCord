package main

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/coder/websocket"
)

type connectedMsg struct {
	conn *websocket.Conn
	ctx  context.Context
}

func (m model) entryUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			name := m.entryInput
			if name == "" {
				return m, nil
			}
			user = name
			m.session = "connecting"
			go connectToServer(program)
			return m, m.spinner.Tick
		case "backspace":
			if len(m.entryInput) > 0 {
				m.entryInput = m.entryInput[:len(m.entryInput)-1]
			}
			return m, nil
		default:
			if len(msg.Text) > 0 {
				m.entryInput += msg.Text
			}
			return m, nil
		}

	}

	return m, nil
}

func (m model) connectionUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connectedMsg:
		m.conn = msg.conn
		m.ctx = msg.ctx
		m.chat = init_chat()
		m.session = "room"
		go listenForMessages(program, m.conn, m.ctx)
		sendServer(m.conn, m.ctx)
		return m, m.spinner.Tick
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func connectToServer(p *tea.Program) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for {
		conn, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://localhost:8080/subscribe?username=%s&color=%d", user, color), nil)
		if err == nil {
			p.Send(connectedMsg{conn: conn, ctx: ctx})
		}
		// return tea.Quit()
	}
}

func (m model) renderEntry() tea.View {
	colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", color)))

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")).Render("⚡ Welcome to TermCord")
	prompt := "Enter your name:"
	input := colorStyle.Render(m.entryInput + "█")
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press Enter when done")

	content := fmt.Sprintf("\n\n  %s\n\n  %s\n  %s\n\n  %s\n", title, prompt, input, hint)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
