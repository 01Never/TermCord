package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/coder/websocket"
)

func main() {

	ctx := context.Background()
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/subscribe", nil)
	if err != nil {
		log.Fatalf("Oof: failed to connect to server: %v", err)
	}

	p := tea.NewProgram(initialModel())
	go listenForMessages(p, conn, ctx)
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

func handler(s string) {
	body := strings.NewReader(s)
	resp, err := http.Post("http://localhost:8080/publish", "text/plain", body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
	defer resp.Body.Close()

}

func (m model) View() tea.View {
	return m.renderChatRoom()
}
