package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)


var (
	applicationName = lipgloss.NewStyle()
	controlTool = lipgloss.NewStyle()
	listRow = lipgloss.NewStyle()
)

func (m model) View() string {
	var newline = "\n\n"
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
		s += controlTool.Render("enter - save ; esc - back");
	}

	if m.viewType == listView {
		for i, n := range m.items {
			prefix := " "
			if m.listIndex == i {
				prefix = "->"
			}

			shortbody := strings.ReplaceAll(n.Body, "\n", "")
			if len(shortbody) > 25 {
				shortbody = shortbody[:25]
			}

			s += listRow.Render(prefix) + "[" + n.Title + "]:" + shortbody + newline
		}

		s += controlTool.Render("a - Add Item ; q - Exit")
	}

	return s
}
