package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	titleView
	bodyView
	priorityView
)

type model struct {
	store       *Store
	viewType    uint
	textarea    textarea.Model
	textinput   textinput.Model
	currentItem Item
	items       []Item
	listIndex   int
}

func NewModel(s *Store) model {
	items, err := s.GetItems()
	if err != nil {
		fmt.Println("Could not create a new Model: %v", err)
		os.Exit(1)
	}

	return model{
		store:     s,
		viewType:  listView,
		textarea:  textarea.New(),
		textinput: textinput.New(),
		items:     items,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Used for batching
	var (
		cmds []tea.Cmd
		cmd tea.Cmd
	)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.viewType {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "a":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.currentItem = Item{}
				m.viewType = titleView
			case "down", "j":
				if m.listIndex <= len(m.items) - 1 {
					m.listIndex++
				}
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
					m.listIndex--
				}
			}
		case titleView:
			switch key {
			case "enter":
				title := m.textinput.Value()
				if title != "" {
					m.currentItem.Title = title
					m.viewType = bodyView
					m.textarea.SetValue("")
					m.textarea.Focus()
					m.textarea.CursorEnd()
				}
			}
		case bodyView:
			switch key {
			case "enter":
				body := m.textarea.Value()
				if body != "" {
					m.currentItem.Body = body
					m.viewType = priorityView
					m.textinput.SetValue("")
					m.textinput.Focus()
					// Empty for Priority
					m.textinput.SetValue("")
				}
			}
		case priorityView:
			switch key {
			case "enter", "ctrl + s":
				priority := m.textinput.Value()
				cint, err := strconv.ParseInt(priority, 10, 64)

				if err != nil {
					fmt.Println("Unable to store priority: ", err)
					os.Exit(1)
				}

				m.currentItem.Priority = cint

				var errr error
				if errr = m.store.CreateItem(m.currentItem); errr != nil {
					log.Fatalf("Could not create item %v", errr)
				}

				m.items, err = m.store.GetItems()
				if err != nil {
					log.Fatalf("Could not get items: %v", err)
				}

				m.viewType = listView
			}
		}
	}

	return m, tea.Batch(cmds...)
}
