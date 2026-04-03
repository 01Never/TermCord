package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Packet struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type serverMsg struct {
	UserID    string `json:"user_id"`
	Context   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type heartBeat struct {
	HeartBeat string `json:"heartBeat"`
}

type model struct {
	chat    chat_model
	session string
	users   []string
	conn    *websocket.Conn
	ctx     context.Context
}

func initialModel(conn *websocket.Conn, ctx context.Context) model {
	return model{
		chat:    init_chat(),
		session: "room",
		conn:    conn,
		ctx:     ctx}
}

// TODO MC periodic stuff for server. Maybe this should go in function that starts
// periodic stuff for the server in general. inside  spilt be sever, client, background etc
func sendServer(conn *websocket.Conn, ctx context.Context) {
	go func() {
		for range time.Tick(10000 * time.Millisecond) {
			data := heartBeat{HeartBeat: "heartbeat: " + user}
			bytes, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("building heartbeat msg error")
			}

			packet := Packet{Type: "heartbeat", Data: bytes}
			err = wsjson.Write(ctx, conn, packet)
			if err != nil {
				fmt.Printf("sending heartbeat error")
			}
		}
	}()
}

func listenForMessages(p *tea.Program, conn *websocket.Conn, ctx context.Context) {
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			fmt.Printf("Something went wrong while cooking")
			return
		}

		var msg serverMsg
		err = json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Printf("JSON is not feeling okay")
		}

		//TODO data needs to include userID. to avoid printing my own message
		p.Send(msg) //this converting data to string then to serverMsg. so bubble.tea understands what this is for the case statments

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

			// convert data to json and insert user name
			msg := serverMsg{UserID: user, Context: string(m.chat.textarea.Value()), Timestamp: 0}
			bytes, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("JSON is not feeling okay")
			}
			// handler(bytes)

			packet := Packet{Type: "msg", Data: bytes}
			err = wsjson.Write(m.ctx, m.conn, packet)
			if err != nil {
				fmt.Printf("Sending via websocket went wrong")
			}

			return m, nil
		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.chat.textarea, cmd = m.chat.textarea.Update(msg)
			return m, cmd
		}
	case serverMsg:
		m.chat.messages = append(m.chat.messages, m.chat.senderStyle.Render(msg.UserID+":")+msg.Context)
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
