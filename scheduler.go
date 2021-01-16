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

// CallListeners calls the event listener of provided event
func (s Scheduler) CallListeners(event Event) {
	eventFn, ok := s.listeners[event.Name]
	if ok {
		go eventFn(event.Payload)
		s.ClearEvent(event)
	} else {
		log.Print("💀 error: couldn't find event listeners attached to ", event.Name)
	}

}

// ClearEvent clears job from database
func (s Scheduler) ClearEvent(event Event) {
	_, err := s.db.Exec(`DELETE FROM "public"."jobs" WHERE "id" = $1`, event.ID)
	if err != nil {
		log.Print("💀 error: ", err)
	}
}

// CheckDueEvents checks and returns due events
func (s Scheduler) CheckDueEvents() []Event {
	events := []Event{}
	rows, err := s.db.Query(`SELECT id, name, payload FROM jobs WHERE runAt < $1`, time.Now())
	if err != nil {
		log.Print("💀 error: ", err)
		return nil
	}
	for rows.Next() {
		evt := Event{}
		rows.Scan(&evt.ID, &evt.Name, &evt.Payload)
		events = append(events, evt)
	}
	return events
}

// Schedule sechedules the provided events
func (s Scheduler) Schedule(event string, payload string, runAt time.Time) {
	_, err := s.db.Exec(`INSERT INTO "public"."jobs" ("name", "payload", "runAt") VALUES ($1, $2, $3)`, event, payload, runAt)
	if err != nil {
		log.Print("schedule insert error: ", err)
	}
}
