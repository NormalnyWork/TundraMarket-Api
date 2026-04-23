package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations done, starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
