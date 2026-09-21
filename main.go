package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seeu359/go-project-278/links"
	"github.com/seeu359/go-project-278/router"
)

func main() {

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("error")
	}
	defer pool.Close()

	q := links.New(pool)
	r := router.New(q)
	r.Run(":8080")
}
