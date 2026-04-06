package main

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	applicationName = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F5F2F2")).
			Background(lipgloss.Color("#2B2A2A")).
			PaddingLeft(16).
			PaddingRight(16)
	controlTool = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A7ACD")).
			Background(lipgloss.Color("#2B2A2A")).
			PaddingLeft(12).
			PaddingRight(12)
	
	listPointer = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A7ACD")).
			PaddingLeft(1).
			PaddingRight(1)
	listTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FEB05D"))
	listDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5F2F2"))
	listPrio = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5F2F2"))
	
	
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

			priority := n.Priority

			s += listPrio.Render("(" + strconv.FormatInt(priority, 10) + ") ") + listPointer.Render(prefix) + listTitle.Render("[" + n.Title + "]: ") + listDesc.Render(shortbody)  + newline
		}

		s += controlTool.Render("a - Add Item ; q - Exit")
	}

	return s
}
