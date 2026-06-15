package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/todo?sslmode=disable"
	}

	db, err := openDB(dsn, 30, time.Second)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	store := NewStore(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := store.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	addr := ":8080"
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, NewHandler(store).Routes()); err != nil {
		log.Fatal(err)
	}
}

// openDB は DB が起動しきるまで ping をリトライする。
func openDB(dsn string, attempts int, wait time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	for i := 0; i < attempts; i++ {
		if err = db.Ping(); err == nil {
			return db, nil
		}
		log.Printf("waiting for db (%d/%d): %v", i+1, attempts, err)
		time.Sleep(wait)
	}
	return nil, err
}
