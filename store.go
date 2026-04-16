package main

import (
	"database/sql"
	"log"
	"time"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Item struct {
	ID           int32
	Title        string
	Body         string
	Priority     int64
	Done         bool
	DateCreate   int32
	DateComplete int32
}

type Store struct {
	conn *sql.DB
}

//go:embed migrations/*.sql
var fs embed.FS

func (s *Store) Init(dbpath string) error {
	var err error

	s.conn, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}

	driver, err := sqlite3.WithInstance(s.conn, &sqlite3.Config{})
	if err != nil {
		return err
	}

	migdifs, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		migdifs,
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
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
			priority,
			done,
			date_create,
			date_complete
		FROM
			items
		ORDER BY priority asc;
	`

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
			&item.Done,
			&item.DateCreate,
			&item.DateComplete,
		)
		items = append(items, item)
	}

	return items, nil
}

func (s *Store) GetPendingItems() ([]Item, error) {
	queryGetItems := `
		SELECT
			id,
			title,
			body,
			priority,
			done,
			date_create,
			date_complete
		FROM
			items
		WHERE done = false
		ORDER BY priority asc;
	`

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
			&item.Done,
			&item.DateCreate,
			&item.DateComplete,
		)
		items = append(items, item)
	}

	return items, nil
}

func (s *Store) GetDoneItems() ([]Item, error) {
	queryGetItems := `
		SELECT
			id,
			title,
			body,
			priority,
			done,
			date_create,
			date_complete
		FROM
			items
		WHERE done = true
		ORDER BY priority asc;
	`

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
			&item.Done,
			&item.DateCreate,
			&item.DateComplete,
		)
		items = append(items, item)
	}

	return items, nil
}

func (s *Store) CreateItem(item Item) error {
	if item.ID == 0 {
		item.ID = int32(time.Now().UTC().Unix())
	}

	item.DateCreate = int32(time.Now().UTC().Unix())

	queryCreateItem := `
		INSERT INTO items (id, title, body, priority, done, date_create, date_complete)
		VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE
		SET 
			title=excluded.title, 
			body=excluded.body, 
			priority=excluded.priority, 
			done=excluded.done, 
			date_create=excluded.date_create, 
			date_complete=excluded.date_complete;
	`

	if _, err := s.conn.Exec(
		queryCreateItem,
		item.ID,
		item.Title,
		item.Body,
		item.Priority,
		false,
		item.DateCreate,
		0,
	); err != nil {
		return err
	}

	return nil
}

func (s *Store) MarkDone(item Item) error {
	if item.ID == 0 {
		log.Fatal("Could not update Item, item does not exit")
	}

	item.DateComplete = int32(time.Now().UTC().Unix())

	queryMarkDone := `
		UPDATE
			items
		SET 
			done = true, 
			date_complete = ?
		WHERE 
			id = ?;
	`
	if _, err := s.conn.Exec(queryMarkDone, item.DateComplete, item.ID); err != nil {
		return err
	}

	return nil
}
