package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	"github.com/dipeshdulal/event-scheduling/customevents"
)

// Event structure
type Event struct {
	ID      uint
	Name    string
	Payload string
	Locked  string
}

var eventListeners = map[string]func(d interface{}){
	"SendEmail": customevents.SendEmail,
	"PayBills":  customevents.PayBills,
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file ", err)
	}

	db := initDBConnection()
	seedDB(db)

	go func() {
		for range interrupt {
			log.Println("\n❌ Interrupt received closing...")
			cancel()
		}
	}()

	<-ctx.Done()
}
