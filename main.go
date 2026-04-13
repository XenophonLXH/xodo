package main

import (
	"fmt"
	"os"
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func getDatabaseDir(dbname string) string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		base = getUserHome()
	}
	xodopath := base + "/xodo/databases/"

	if err := os.MkdirAll(xodopath, 0755); err != nil {
		log.Fatalf("Could not find or create database directory: %v", err)
	}

	return xodopath + dbname
}

func getUserHome() string {
	if homedir, err := os.UserHomeDir(); err != nil {
		log.Fatalf("Could not determine user's home directory %v", err)
	} else {
		return homedir
	}
	return ""
}

func main() {

	// Default to xodo as db
	// if none is passed
	var dbname = "xodo"
	if len(os.Args) > 1 {
		dbname = os.Args[1]
	}
	dbpath := getDatabaseDir(dbname)

	// Get connected to the database
	s :=  new(Store)

	if err := s.Init(dbpath); err != nil {
		fmt.Println("Could not initiate store:", err)
		os.Exit(1)
	}

	// Initiate the model
	m := NewModel(s)

	// Set list title
	m.listName = strings.ToUpper(dbname)

	// Bubble tea program
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Unable to start TUI %v", err)
	}

}
