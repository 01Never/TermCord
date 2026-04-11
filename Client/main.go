package main

import (
	"context"
	"encoding/json"
	"example/TermCord/shared"
	"fmt"
	"math/rand"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/coder/websocket"
)

type model struct {
	chat       chat_model
	session    string
	users      []string
	conn       *websocket.Conn
	ctx        context.Context
	entryInput string
}

var user string = fmt.Sprintf("USER_%03d", rand.Intn(1000))
var color int = rand.Intn(256)
var program *tea.Program

func main() {
	p := tea.NewProgram(initialModel())
	program = p

	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.session {
	case "entry":
		return m.entryUpdate(msg)
	case "room":
		return m.roomUpdate(msg)
	}
	return m, nil
}

func (m model) View() tea.View {
	switch m.session {
	case "entry":
		return m.renderEntry()
	case "room":
		return m.renderChatRoom()
	}
	return m.renderEntry()
}

func listenForMessages(p *tea.Program, conn *websocket.Conn, ctx context.Context) {
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			fmt.Printf("Something went wrong while cooking")
			return
		}

		var packet shared.Packet
		err = json.Unmarshal(data, &packet)
		if err != nil {
			fmt.Printf("JSON inside a JSON!!")
		}

		switch packet.Type {
		case "PostMsg":
			var msg shared.PostMsg
			err = json.Unmarshal(packet.Data, &msg)
			if err != nil {
				fmt.Printf("Failed to unmarshal MsgPosted")
			}
			p.Send(msg)

		case "UserJoined":
			var msg shared.UserJoined
			err = json.Unmarshal(packet.Data, &msg)
			if err != nil {
				fmt.Printf("Failed to unmarshal UserJoined")
			}
			p.Send(msg)

		case "UserLeft":
			var msg shared.UserLeft
			err = json.Unmarshal(packet.Data, &msg)
			if err != nil {
				fmt.Printf("Failed to unmarshal UserLeft")
			}
			p.Send(msg)

		case "RoomState":
			var msg shared.RoomState
			err = json.Unmarshal(packet.Data, &msg)
			if err != nil {
				fmt.Printf("Failed to unmarshal RoomState")
			}
			p.Send(msg)
		}
	}
}
