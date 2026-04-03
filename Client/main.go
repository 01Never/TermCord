package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/coder/websocket"
)

var user string = fmt.Sprintf("USER_%03d", rand.Intn(1000))

func main() {

	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/subscribe?username="+user, nil)
	if err != nil {
		log.Fatalf("Oof: failed to connect to server: %v", err)
	}

	p := tea.NewProgram(initialModel(conn, ctx))
	go listenForMessages(p, conn, ctx)
	sendServer(conn, ctx)

	_, err = p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.session {
	case "room":
		return m.roomUpdate(msg)

	}

	return m, nil
}

func handler(b []byte) {
	body := bytes.NewReader(b)
	resp, err := http.Post("http://localhost:8080/publish", "application/json", body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
	defer resp.Body.Close()

}

func (m model) View() tea.View {
	return m.renderChatRoom()
}
