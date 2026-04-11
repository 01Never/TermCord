package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type SidebarItem struct {
	Name  string
	Color int
}

var activeChannelStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("255")).
	Bold(true).
	Padding(0, 1)

func (m model) Init() tea.Cmd {
	return nil
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

	//currentCategory := ""
	for _, ch := range m.chat.onlineUsers {
		// if ch.Category != currentCategory {
		// 	currentCategory = ch.Category
		// 	b.WriteString(categoryStyle.Width(sidebarW).Render(ch.Category) + "\n")
		// 	usedLines++
		// }
		nameStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("%d", ch.Color))).
			Padding(0, 1)
		b.WriteString(nameStyle.Width(sidebarW).Render(ch.Name) + "\n")
		usedLines++
	}

	// Fill remaining lines so the right border extends to the bottom
	for i := 0; i < (m.chat.viewport.Height()+m.chat.textarea.Height())-usedLines; i++ {
		b.WriteString(strings.Repeat(" ", sidebarW) + "\n")
	}

	// Wrap the whole sidebar in the right-border style
	return b.String()
}
