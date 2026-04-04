package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID int32
	Title string
	Body string
	Priority int32
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error

	s.conn, err = sql.Open("sqlite3", "./xodo.db")

	if err != nil {
		return err
	}

	queryTableCreate := `
		CREATE TABLE IF NOT EXISTS items (
			id integer not null primary key
			title text not null
			body text
			priority integer not null
		);
	`

	if _, err := s.conn.Exec(queryTableCreate); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetItems() ([]Item, error) {
	return nil, nil
}

func (s *Store) CreateItem(item Item) error {
	return nil
}
