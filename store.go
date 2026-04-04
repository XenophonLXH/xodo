package main

import (
	"database/sql"
	"time"
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
	queryGetItems := `
		SELECT
			id,
			title,
			body,
			priority
		FROM
			items
		ORDER priortiy asc;
	`;

	rows, err := s.conn.Query(queryGetItems)

	if err != nil {
		return nil, err
	}

	items := []Item{}
	defer rows.Close()
	for rows.Next() {
		item := Item{}
		rows.Scan(
			&item.ID,
			&item.Title,
			&item.Body,
			&item.Priority,
		)
		items = append(items, item)
	}

	return nil, nil
}

func (s *Store) CreateItem(item Item) error {
	if item.ID == 0 {
		item.ID = int32(time.Now().UTC().Unix())
	}

	queryCreateItem := `
		INSERT INTO items (id, title, body, priority)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE
		SET title=exluded.title, body=exluded.body, priority=exluded.priority;
	`

	if _, err := s.conn.Exec(queryCreateItem, item.ID, item.Title, item.Body, item.Priority); err != nil {
		return err
	}

	return nil
}
