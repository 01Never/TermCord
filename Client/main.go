package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Title bar — bold white, border only on the bottom to act as a divider
var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("255")).
	Bold(true).
	Padding(0, 1).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(lipgloss.Color("240"))

	// Category headers — dimmer than channel names
var categoryStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("238")).
	Bold(true).
	Padding(0, 1)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

type Channel struct {
	Name     string
	Category string
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	channels    []Channel
	err         error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "┃"
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	// Remove cursor line styling
	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		channels:    []Channel{},
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func handler(s string) {
	body := strings.NewReader(s)
	resp, err := http.Post("http://localhost:8080/helloWorld", "text/plain", body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height())

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case "enter":
			//TODO here to place user name from server.
			//TODO here to send message to server
			handler(m.senderStyle.Render("You: ") + m.textarea.Value())

			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, nil
		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() tea.View {
	viewportView := m.viewport.View()
	var content string
	content = lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderSidebar(),
		(viewportView + "\n" + m.textarea.View()),
	)

	v := tea.NewView(content)
	c := m.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}

func (m model) renderSidebar() string {
	const sidebarW = 22
	var b strings.Builder

	b.WriteString(titleStyle.Width(sidebarW).Render("⚡ TermCord") + "\n")
	usedLines := 2

	// currentCategory := ""
	// for ch := range m.channels {
	// 	if ch.Category != currentCategory {
	// 		currentCategory = ch.Category
	// 		b.WriteString(categoryStyle.Width(sidebarW).Render(ch.Category) + "\n")
	// 		usedLines++
	// 	}
	// 	label := "# " + ch.Name
	// 	if i == m.activeChannel {
	// 		b.WriteString(activeChannelStyle.Width(sidebarW).Render(label) + "\n")
	// 	} else {
	// 		b.WriteString(channelStyle.Width(sidebarW).Render(label) + "\n")
	// 	}
	// 	usedLines++
	//}

	// Fill remaining lines so the right border extends to the bottom
	for i := 0; i < (m.viewport.Height()+m.textarea.Height())-usedLines; i++ {
		b.WriteString(strings.Repeat(" ", sidebarW) + "\n")
	}

	// Wrap the whole sidebar in the right-border style
	return b.String()
}
