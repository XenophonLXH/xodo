package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

    "charm.land/bubbles/v2/textarea"
    "charm.land/bubbles/v2/textinput"
    "charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

const (
	listView uint = iota
	titleView
	bodyView
	priorityView
)

type listMode int

const (
	pending = iota
	done
	all
)

type model struct {
	listName	string
	store       *Store
	viewport	viewport.Model
	viewType    uint
	textarea    textarea.Model
	textinput   textinput.Model
	currentItem Item
	items       []Item
	listIndex   int
	listMode	listMode
}

func NewModel(s *Store) model {
	items, err := s.GetPendingItems()
	if err != nil {
		fmt.Printf("Could not create a new Model: %v", err)
		os.Exit(1)
	}

	viewp := viewport.New(viewport.WithWidth(5), viewport.WithHeight(5))

	return model{
		store:     s,
		viewport: viewp,
		viewType:  listView,
		textarea:  textarea.New(),
		textinput: textinput.New(),
		items:     items,
		listMode:	0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func handleListMode(m model) model {
	var err error
	switch m.listMode {
	case 0:
		m.items, err = m.store.GetPendingItems()
	case 1:
		m.items, err = m.store.GetDoneItems()
	case 2:
		m.items, err = m.store.GetItems()
	}
	if err != nil {
		log.Fatalf("Could not fetch items: %v", err)
	}
	return m
}

func handleListView(
	m model,
	key string,
) model {
	switch key {
	case "tab":
		m.listIndex = 0
		switch m.listMode {
		case 0:
			m.listMode = 1

			var err error
			m.items, err = m.store.GetDoneItems()
			if err != nil {
				log.Fatalf("Could not get done items: %v", err)
			}

		case 1:
			m.listMode = 2

			var err error
			m.items, err = m.store.GetItems()
			if err != nil {
				log.Fatalf("Could not get done items: %v", err)
			}

		case 2:
			m.listMode = 0

			var err error
			m.items, err = m.store.GetPendingItems()
			if err != nil {
				log.Fatalf("Could not get done items: %v", err)
			}

		}

	case "a":
		m.textinput.SetValue("")
		m.textinput.Focus()
		m.currentItem = Item{}
		m.viewType = titleView

	case "down", "j":
		if m.listIndex < len(m.items)-1 {
			m.listIndex++
		} else {
			m.listIndex = 0
		}

	case "up", "k":
		if m.listIndex > 0 {
			m.listIndex--
		} else {
			m.listIndex = len(m.items) - 1
		}

	case "d":
		m.currentItem = m.items[m.listIndex]

		if !m.currentItem.Done {
			var err error
			if err = m.store.MarkDone(m.currentItem); err != nil {
				log.Fatalf("Could not mark the item as done %v", err)
			}

			m.listIndex = 0
			m = handleListMode(m)

		}

	case "p":
		m.currentItem = m.items[m.listIndex]

		if m.currentItem.Done {
			var err error
			if err = m.store.MarkPending(m.currentItem); err != nil {
				log.Fatalf("Could not mark the item as done %v", err)
			}

			m.listIndex = 0
			
			m = handleListMode(m)
		}


	case "i":
		m.currentItem = m.items[m.listIndex]
		m.viewType = bodyView
		m.textarea.SetValue(m.currentItem.Body)
		m.textarea.Focus()
		m.textarea.CursorEnd()
	}
	return m
}

func handleTitleView(m model, key string) model {
	switch key {
	case "enter", "ctrl+s":
		title := m.textinput.Value()
		if title != "" {
			m.currentItem.Title = title
			m.viewType = bodyView
			m.textarea.SetValue("")
			m.textarea.Focus()
			m.textarea.CursorEnd()
		}

	case "q":
		m.viewType = listView
	}
	return m
}

func handleBodyView(m model, key string) model {
	switch key {
	case "enter":
		body := m.textarea.Value()
		if body != "" {
			m.currentItem.Body = body
			m.viewType = priorityView
			m.textinput.SetValue("")
			m.textinput.Focus()
			// Empty for Priority
			if m.currentItem.Priority != 0 {
				m.textinput.SetValue(strconv.FormatInt(m.currentItem.Priority, 10))
			} else {
				m.textinput.SetValue("")
			}
		}
	case "esc":
		m.viewType = listView
	}
	return m
}

func handlePriorityView(m model, key string) model {
	switch key {
	case "enter", "ctrl+s":
		priority := m.textinput.Value()
		cint, err := strconv.ParseInt(priority, 10, 64)
		if err != nil {
			fmt.Println("Unable to store priority: ", err)
			os.Exit(1)
		}

		m.currentItem.Priority = cint

		var errr error
		if errr = m.store.CreateItem(m.currentItem); errr != nil {
			log.Fatalf("Could not create item: %v", errr)
		}

		m.listIndex = 0
		m.listMode = 0
		m.items, err = m.store.GetPendingItems()
		if err != nil {
			log.Fatalf("Could not get items: %v", err)
		}

		m.viewType = listView
	}
	return m
}

// sets the viewport width and height
// when the terminal resizes
func handleWindowResize(m model, msg tea.WindowSizeMsg) model {
	m.viewport.SetWidth(msg.Width)
	m.viewport.SetHeight(msg.Height - 10)
	return m
}

func mustExit(key string) bool {
	mustExit := false
	switch key {
		case "ctrl+w":
			mustExit = true
	}
	return mustExit
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Used for batching
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()

		mustExit := mustExit(key)
		if mustExit {
			return m, tea.Quit
		}

		switch m.viewType {

		case listView:
			m = handleListView(
				m,
				key,
			)

		case titleView:
			m = handleTitleView(m, key)

		case bodyView:
			m = handleBodyView(m, key)

		case priorityView:
			m = handlePriorityView(m, key)

		}

		case tea.WindowSizeMsg:
			m = handleWindowResize(m, msg)
	}

	return m, tea.Batch(cmds...)
}
