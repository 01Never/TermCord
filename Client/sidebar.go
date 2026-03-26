package main

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Channel struct {
	Name     string
	Category string
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// Styling
var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("255")).
	Bold(true).
	Padding(0, 1).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(lipgloss.Color("240"))

var categoryStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("238")).
	Bold(true).
	Padding(0, 1)

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
	for i := 0; i < (m.chat.viewport.Height()+m.chat.textarea.Height())-usedLines; i++ {
		b.WriteString(strings.Repeat(" ", sidebarW) + "\n")
	}

	// Wrap the whole sidebar in the right-border style
	return b.String()
}
