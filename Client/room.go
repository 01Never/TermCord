package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"example/TermCord/shared"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func initialModel() model {
	return model{
		session: "entry",
	}
}

// TODO MC periodic stuff for server. Maybe this should go in function that starts
// periodic stuff for the server in general. inside  spilt be sever, client, background etc
func sendServer(conn *websocket.Conn, ctx context.Context) {
	go func() {
		for range time.Tick(10000 * time.Millisecond) {
			data := shared.HeartBeat{HeartBeat: "heartbeat: " + user}
			bytes, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("building heartbeat msg error")
			}

			packet := shared.Packet{Type: "heartbeat", Data: bytes}
			err = wsjson.Write(ctx, conn, packet)
			if err != nil {
				fmt.Printf("sending heartbeat error")
			}
		}
	}()
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
			if string(m.chat.textarea.Value()) == "" {
				return m, nil
			}

			msg := shared.PostMsg{UserId: user, Msg: string(m.chat.textarea.Value()), Color: color}
			bytes, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Failed to Marshal PostMsg")
			}

			packet := shared.Packet{Type: "PostMsg", Data: bytes}
			err = wsjson.Write(m.ctx, m.conn, packet)
			if err != nil {
				fmt.Printf("Sending Packet went wrong")
			}

			m.chat.messages = append(m.chat.messages, m.chat.senderStyle.Render(user+":")+m.chat.textarea.Value())
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
	case shared.PostMsg:
		m.chat.messages = append(m.chat.messages, lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", msg.Color))).Render(msg.UserId+":")+msg.Msg)
		m.chat.viewport.SetContent(lipgloss.NewStyle().Width(m.chat.viewport.Width()).Render(strings.Join(m.chat.messages, "\n")))
		m.chat.textarea.Reset()
		m.chat.viewport.GotoBottom()

		return m, nil

	case shared.UserJoined:
		//something
		onlineUser := SidebarItem{Name: msg.UserID, Category: "Development"}
		m.chat.onlineUsers = append(m.chat.onlineUsers, onlineUser)
		return m, nil

	case shared.RoomState:
		m.chat.onlineUsers = nil
		for _, userID := range msg.OnlineUsers {
			m.chat.onlineUsers = append(m.chat.onlineUsers, SidebarItem{Name: userID, Category: "Development"})
		}
		return m, nil

	case shared.UserLeft:
		for i, u := range m.chat.onlineUsers {
			if u.Name == msg.UserID {
				m.chat.onlineUsers = append(m.chat.onlineUsers[:i], m.chat.onlineUsers[i+1:]...)
				break
			}
		}
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
