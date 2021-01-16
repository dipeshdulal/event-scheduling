package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func initDBConnection() *sql.DB {
	connStr := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Panic("couldn't connect to database", err)
	}

	return db
}

func seedDB(db *sql.DB) error {
	log.Print("💾 Seeding database with table...")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "public"."jobs" (
			"id"      SERIAL PRIMARY KEY,
			"name"    varchar(50) NOT NULL,
			"payload" text,
			"runAt"   TIMESTAMP NOT NULL
		)
	`)
	return err
}
