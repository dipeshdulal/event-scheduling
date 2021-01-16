package main

import (
	"database/sql"
	"log"
)

func seedDB(db *sql.DB) error {
	log.Print("💾 Seeding database with table...")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "public"."test" (
			"id" integer NOT NULL,
			"name" text NOT NULL
		)
	`)
	return err
}
