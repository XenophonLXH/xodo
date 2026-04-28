package main

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	titleFG = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F5F2F2"))

	titleBG = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FEB05D")).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center)

	listNameFG = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FEB05D")).
			Padding(0, 1)

	controlTool = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A7ACD")).
			Margin(0).
			Align(lipgloss.Center)

	listPointer = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5A7ACD")).
			Padding(0, 1)
	listTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FEB05D"))
	listDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5F2F2"))
	listPrio = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F5F2F2"))

	viewModeInactive = lipgloss.NewStyle().
				Background(lipgloss.Color("#2B2A2A")).
				Padding(0, 1).
				Margin(0, 1).
				Align(lipgloss.Center)

	viewModeActive = lipgloss.NewStyle().
			Background(lipgloss.Color("#5A7ACD")).
			Foreground(lipgloss.White).
			Padding(0, 1).
			Margin(0, 1).
			Align(lipgloss.Center)
)

func (m model) View() tea.View {
	newline := "\n\n"
	// Title
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 80
	}
	titleBG = titleBG.Width(termWidth - 2)
	listName := listNameFG.Render(m.listName)
	title := titleFG.Render("Your") + listName + titleFG.Render("TODO list!")
	s := titleBG.Render(title) + "\n"

	// Current View Type
	s += renderListMode(m.listMode, termWidth)

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

	var viewPortContent string
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
			viewPortContent += listPrio.Render(
				"("+strconv.FormatInt(priority, 10)+") ",
			) + listPointer.Render(prefix) + listTitle.Render("["+n.Title+"]: ") + listDesc.Render(shortbody) + newline
		}
	}

	// Viewport 
	m.viewport.SetContent(viewPortContent)

	// Handle viewport scrolling
	if m.viewType == listView {
		const lineHeight = 2
		currItemOffset := m.listIndex * lineHeight
		viewPortHeight := m.viewport.Height()

		if currItemOffset < m.viewport.YOffset() {
			m.viewport.SetYOffset(currItemOffset)
		} else if currItemOffset + lineHeight > m.viewport.YOffset() + viewPortHeight {
			m.viewport.SetYOffset(currItemOffset + lineHeight - viewPortHeight)
		}
	}

	s += m.viewport.View()

	// Help tool
	s += newline
	s += renderHelpTool(m, termWidth)

	v := tea.NewView(s)
	v.AltScreen = true

	return v
}

func renderListMode(lm listMode, t int) (s string) {
	newline := "\n\n"

	viewModeActive = viewModeActive.Width((t / 3) - 2)
	viewModeInactive = viewModeInactive.Width((t / 3) - 2)

	if lm == 0 {
		return  viewModeActive.Render("Pending") + viewModeInactive.Render("Done") + viewModeInactive.Render("All") + newline
	}

	if lm == 1 {
		return  viewModeInactive.Render("Pending") + viewModeActive.Render("Done") + viewModeInactive.Render("All") + newline
	}

	if lm == 2 {
		return  viewModeInactive.Render("Pending") + viewModeInactive.Render("Done") + viewModeActive.Render("All") + newline
	}

	return ""
}

func renderHelpTool(m model, t int) (s string) {
	controlTool = controlTool.Width(t)
	if len(m.items) == 0 && (m.listMode == 0 || m.listMode == 2) {
		return controlTool.Render("a - add ; q - quit ; i - edit ; d - done")
	}

	if len(m.items) == 0 && m.listMode == 1 {
		return controlTool.Render("a - add ; q - quit ; i - edit ; p - pending")
	}

	if m.items[m.listIndex].Done {
		return controlTool.Render("a - add ; q - quit ; i - edit ; p - pending")
	} else {
		return controlTool.Render("a - add ; q - quit ; i - edit ; d - done")
	}

}
