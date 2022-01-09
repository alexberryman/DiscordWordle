package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"log"
	"os"
)

func main() {
	migrations := &migrate.FileMigrationSource{
		Dir: "internal/turnips/schema",
	}

	DatabaseUrl := os.Getenv("DATABASE_URL")
	if DatabaseUrl == "" {
		log.Fatal("databaseUrl must be set")
	}

	db, err := sql.Open("postgres", DatabaseUrl)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal("migrations failed to up", err)
	}

	log.Printf("Applied %d migrations!\n", n)
}
