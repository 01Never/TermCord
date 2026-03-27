package main

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	chat chat_model
}

func initialModel() model {

	return model{
		chat: init_chat(),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// this should go under chat.go
		m.chat.viewport.SetWidth(msg.Width)
		m.chat.textarea.SetWidth(msg.Width)
		m.chat.viewport.SetHeight(msg.Height - m.chat.textarea.Height() - 2) // todo this 2 accounts for the header. this can be done better

		if len(m.chat.messages) > 0 {
			// Wrap content before setting it.
			m.chat.viewport.SetContent(lipgloss.NewStyle().Width(m.chat.viewport.Width()).Render(strings.Join(m.chat.messages, "\n")))
		}
		m.chat.viewport.GotoBottom()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println(m.chat.textarea.Value())
			return m, tea.Quit
		case "enter":
			//TODO here to place user name from server.
			//TODO here to send message to server
			handler(m.chat.senderStyle.Render("You: ") + m.chat.textarea.Value())

			m.chat.messages = append(m.chat.messages, m.chat.senderStyle.Render("You: ")+m.chat.textarea.Value())
			m.chat.viewport.SetContent(lipgloss.NewStyle().Width(m.chat.viewport.Width()).Render(strings.Join(m.chat.messages, "\n")))
			m.chat.textarea.Reset()
			m.chat.viewport.GotoBottom()
			return m, nil
		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.chat.textarea, cmd = m.chat.textarea.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.chat.textarea, cmd = m.chat.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) renderChatRoom() tea.View {
	sidebarView := m.renderSidebar()
	chatView := m.renderChat()

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, chatView)

	v := tea.NewView(content)
	c := m.chat.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(m.chatViewportPanel())
		c.X += lipgloss.Width(sidebarView)
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}
