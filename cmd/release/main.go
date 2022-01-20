package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	migrate "github.com/rubenv/sql-migrate"
	"os"
)

func main() {
	migrations := &migrate.FileMigrationSource{
		Dir: "internal/wordle/schema",
	}

	DatabaseUrl := os.Getenv("DATABASE_URL")
	if DatabaseUrl == "" {
		log.Fatal().Msg("databaseUrl must be set")
	}

	db, err := sql.Open("postgres", DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal().Err(err).Msg("migrations failed to up")
	}

	log.Info().Msgf("Applied %d migrations!", n)
}
