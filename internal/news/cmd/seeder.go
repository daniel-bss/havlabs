package main

import (
	"context"

	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func runSeeder(connPool *pgxpool.Pool) {
	store := db.NewStore(connPool)

	store.CreateNews(context.Background(), db.CreateNewsParams{CreatorUsername: "aaa", Title: "bbb", Content: "ccc"})
}
