package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

func handler(s string) {
	body := strings.NewReader(s)
	resp, err := http.Post("http://localhost:8080/helloWorld", "text/plain", body)
	if err != nil {
		panic(err)
		// this should not crash the program
	}
	defer resp.Body.Close()

}

func (m model) View() tea.View {
	return m.render_chat_room()
}
