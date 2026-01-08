package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"server/internal/env"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS env vars")
	}
	cfg := config{
		addr: env.GetEnv("PORT", "8080"),
		db: dbConfig{
			dsn: env.GetEnv("GOOSE_DBSTRING", ""),
		},
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	api := application{
		config: cfg,
		dbConn: conn,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run(api.mount()); err != nil {
		log.Fatal(err)
	}
}
