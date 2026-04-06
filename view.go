package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	applicationName = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F5F2F2")).
			Background(lipgloss.Color("#2B2A2A")).
			PaddingLeft(8).
			PaddingRight(8)
	controlTool = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A7ACD")).
			Background(lipgloss.Color("#2B2A2A")).
			PaddingLeft(4).
			PaddingRight(4)
	
	listRow     = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A7ACD")).
			PaddingLeft(4).
			PaddingRight(4)
)

func (m model) View() string {
	newline := "\n\n"
	s := applicationName.Render("Your TODO List!") + newline

	if m.viewType == titleView {
		s += "Title: " + newline
		s += m.textinput.View() + newline
		s += controlTool.Render("enter - save ; esc - discard")
	}

	if m.viewType == bodyView {
		s += "Description: " + newline
		s += m.textarea.View() + newline
		s += controlTool.Render("enter - save ; esc - back")
	}

	if m.viewType == priorityView {
		s += "Priority: " + newline
		s += m.textinput.View() + newline
		s += controlTool.Render("enter - save ; esc - back")
	}

	if m.viewType == listView {
		for i, n := range m.items {
			prefix := " "
			if m.listIndex == i {
				prefix = "-> "
			}

			shortbody := strings.ReplaceAll(n.Body, "\n", "")
			if len(shortbody) > 25 {
				shortbody = shortbody[:25]
			}

			s += listRow.Render(prefix) + "[" + n.Title + "]: " + shortbody + newline
		}

		s += controlTool.Render("a - Add Item ; q - Exit")
	}

	return s
}
