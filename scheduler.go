package main

import (
	"database/sql"
	"log"
	"time"
)

// Scheduler data structure
type Scheduler struct {
	db        *sql.DB
	listeners Listeners
}

// Listeners has attached event listeners
type Listeners map[string]ListenFunc

// ListenFunc function that listens to events
type ListenFunc func(string)

// Event structure
type Event struct {
	ID      uint
	Name    string
	Payload string
}

// NewScheduler creates a new scheduler
func NewScheduler(db *sql.DB, listeners Listeners) Scheduler {
	return Scheduler{
		db:        db,
		listeners: listeners,
	}
}

// AddListener adds the listener function to Listeners
func (s Scheduler) AddListener(event string, listenFunc ListenFunc) {
	s.listeners[event] = listenFunc
}

// CheckDueEvents checks and returns due events
func (s Scheduler) CheckDueEvents() {
	// var event Event
	// s.db.Query(`SELECT id, name, payload FROM jobs WHERE locked = 0`)
}

// Schedule sechedules the provided events
func (s Scheduler) Schedule(event string, payload string, runAt time.Time) {
	_, err := s.db.Exec(`INSERT INTO "public"."jobs" ("name", "payload", "runAt") VALUES ($1, $2, $3)`, event, payload, runAt)
	if err != nil {
		log.Print("schedule insert error: ", err)
	}
}
