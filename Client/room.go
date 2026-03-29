package main

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/coder/websocket"
)

//var lastmsg string

type serverMsg string

type model struct {
	chat    chat_model
	session string
}

func initialModel() model {
	return model{
		chat:    init_chat(),
		session: "room"}
}

func listenForMessages(p *tea.Program, conn *websocket.Conn, ctx context.Context) {
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			fmt.Printf("Something went wrong while cooking")
			return
		}

		//TODO data needs to include userID. to avoid printing my own message
		p.Send(serverMsg(string(data))) //this converting data to string then to serverMsg. so bubble.tea understands what this is for the case statments

	}
}

func (m model) roomUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			handler(m.chat.textarea.Value())
			//lastmsg = string(m.chat.textarea.Value())

			//TODO instead of updating the message locally just send to the server and when we get the messgae back we prit
			//avoiding the need for json messages yet.

			// m.chat.messages = append(m.chat.messages, m.chat.senderStyle.Render("ME: ")+m.chat.textarea.Value())
			// m.chat.viewport.SetContent(lipgloss.NewStyle().Width(m.chat.viewport.Width()).Render(strings.Join(m.chat.messages, "\n")))
			// m.chat.textarea.Reset()
			// m.chat.viewport.GotoBottom()
			return m, nil
		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.chat.textarea, cmd = m.chat.textarea.Update(msg)
			return m, cmd
		}
	case serverMsg:
		m.chat.messages = append(m.chat.messages, m.chat.senderStyle.Render("ME: ")+string(msg))
		m.chat.viewport.SetContent(lipgloss.NewStyle().Width(m.chat.viewport.Width()).Render(strings.Join(m.chat.messages, "\n")))
		m.chat.textarea.Reset()
		m.chat.viewport.GotoBottom()
		return m, nil

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
