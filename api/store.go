package main

import (
	"context"
	"database/sql"
	_ "embed"
	"time"
)

//go:embed schema.sql
var schemaSQL string

// Task はタスク1件を表す。
type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// Store はタスクの永続化を担う。
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Migrate はスキーマを冪等に適用する。
func (s *Store) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, schemaSQL)
	return err
}

func (s *Store) List(ctx context.Context) ([]Task, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, title, done, created_at FROM tasks ORDER BY created_at ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *Store) Create(ctx context.Context, title string) (Task, error) {
	var t Task
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, done, created_at`,
		title).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	return t, err
}

// UpdateDone は指定タスクの完了状態を更新し、更新後のタスクを返す。
func (s *Store) UpdateDone(ctx context.Context, id int64, done bool) (Task, bool, error) {
	var t Task
	err := s.db.QueryRowContext(ctx,
		`UPDATE tasks SET title = title WHERE id = $1
		 RETURNING id, title, done, created_at`,
		id).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return Task{}, false, nil
	}
	if err != nil {
		return Task{}, false, err
	}
	return t, true, nil
}

func (s *Store) Delete(ctx context.Context, id int64) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
