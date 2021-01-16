package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

// Test structure
type Test struct {
	ID   uint
	Name string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	interrupt := make(chan os.Signal, 1)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file ", err)
	}

	db := initDBConnection()

	rows, err := db.Query("SELECT id, name FROM public.test")

	if err != nil {
		log.Panic(err)
	}

	go func() {
		for rows.Next() {
			var test Test
			rows.Scan(&test.ID, &test.Name)
			log.Printf("test: %v\n", test)
		}
	}()

	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for range interrupt {
			log.Print("Interrupt received closing...")
			cancel()
		}
	}()

	<-ctx.Done()
}
