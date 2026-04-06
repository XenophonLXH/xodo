package main

import (
	"fmt"
	"os"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Get connected to the database
	s :=  new(Store)

	if err := s.Init(); err != nil {
		fmt.Println("Could not initiate store:", err)
		os.Exit(1)
	}

	// Initiate the model
	m := NewModel(s)

	// Bubble tea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Unable to start TUI %v", err)
	}

}
