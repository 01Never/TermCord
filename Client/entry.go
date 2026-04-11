package main

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/textarea"
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
			return m, connectToServer
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

	case connectedMsg:
		m.conn = msg.conn
		m.ctx = msg.ctx
		m.chat = init_chat()
		m.session = "room"
		go listenForMessages(program, m.conn, m.ctx)
		sendServer(m.conn, m.ctx)
		return m, textarea.Blink
	}

	return m, nil
}

func connectToServer() tea.Msg {
	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/subscribe?username="+user, nil)
	if err != nil {
		return tea.Quit()
	}
	return connectedMsg{conn: conn, ctx: ctx}
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
