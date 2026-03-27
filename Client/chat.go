package main

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

type chat_model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	channels    []Channel
	err         error
}

func init_chat() chat_model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "┃"
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome to the chat room! Type a message and press Enter to send.`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return chat_model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		channels:    []Channel{},
		err:         nil,
	}
}

func (m model) chatViewportPanel() string {
	return titleStyle.Width(m.chat.viewport.Width()).Render("#general") + "\n" + m.chat.viewport.View()
}

func (m model) renderChat() string {
	var b strings.Builder
	b.WriteString(m.chatViewportPanel() + "\n" + m.chat.textarea.View())
	return b.String()
}
