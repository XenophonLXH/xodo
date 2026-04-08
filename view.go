package main

import (
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
)

var (
	applicationName = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F5F2F2")).
			Background(lipgloss.Color("#2B2A2A")).
			Padding(0, 35)
	controlTool = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A7ACD")).
			Background(lipgloss.Color("#2B2A2A")).
			Padding(0, 22)

	listPointer = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A7ACD")).
			Padding(0, 1)
	listTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FEB05D"))
	listDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5F2F2"))
	listPrio = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5F2F2"))
)

func (m model) View() tea.View {
	newline := "\n\n"
	s := applicationName.Render("Your TODO List!") + newline

	if m.viewType == titleView {
		s += listTitle.Render("Title: ") + newline
		s += m.textinput.View() + newline
		s += controlTool.Render("enter - save ; esc - discard ; ctrl + w - quit")
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
				shortbody = shortbody[:25] + "...."
			}

			priority := n.Priority

			s += listPrio.Render("("+strconv.FormatInt(priority, 10)+") ") + listPointer.Render(prefix) + listTitle.Render("["+n.Title+"]: ") + listDesc.Render(shortbody) + newline
		}

		s += controlTool.Render("a - add ; q - quit ; i - edit ; d - done")
	}

	return s
}
