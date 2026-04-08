package main

import (
	"fmt"
	"os"
	"log"

	tea "charm.land/bubbletea/v2"
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
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Unable to start TUI %v", err)
	}

}
