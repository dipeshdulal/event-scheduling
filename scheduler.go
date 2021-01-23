package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler data structure
type Scheduler struct {
	db          *sql.DB
	listeners   Listeners
	cron        *cron.Cron
	cronEntries map[string]cron.EntryID
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
	Cron    string
}

// NewScheduler creates a new scheduler
func NewScheduler(db *sql.DB, listeners Listeners) Scheduler {

	return Scheduler{
		db:          db,
		listeners:   listeners,
		cron:        cron.New(),
		cronEntries: map[string]cron.EntryID{},
	}

}

// AddListener adds the listener function to Listeners
func (s Scheduler) AddListener(event string, listenFunc ListenFunc) {
	s.listeners[event] = listenFunc
}

// callListeners calls the event listener of provided event
func (s Scheduler) callListeners(event Event) {
	eventFn, ok := s.listeners[event.Name]
	if ok {
		go eventFn(event.Payload)
		_, err := s.db.Exec(`DELETE FROM "public"."jobs" WHERE "id" = $1`, event.ID)
		if err != nil {
			log.Print("💀 error: ", err)
		}
	} else {
		log.Print("💀 error: couldn't find event listeners attached to ", event.Name)
	}

}

// CheckEventsInInterval checks the event in given interval
func (s Scheduler) CheckEventsInInterval(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("⏰ Ticks Received...")
				events := s.checkDueEvents()
				for _, e := range events {
					s.callListeners(e)
				}
			}

		}
	}()
}

// checkDueEvents checks and returns due events
func (s Scheduler) checkDueEvents() []Event {
	events := []Event{}
	rows, err := s.db.Query(`SELECT "id", "name", "payload" FROM "public"."jobs" WHERE "runAt" < $1 AND "cron"='-'`, time.Now())
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
	log.Print("🚀 Scheduling event ", event, " to run at ", runAt)
	_, err := s.db.Exec(`INSERT INTO "public"."jobs" ("name", "payload", "runAt") VALUES ($1, $2, $3)`, event, payload, runAt)
	if err != nil {
		log.Print("schedule insert error: ", err)
	}
}

// ScheduleCron schedules a cron job
func (s Scheduler) ScheduleCron(event string, payload string, cron string) {
	log.Print("🚀 Scheduling event ", event, " with cron string ", cron)
	entryID, ok := s.cronEntries[event]
	if ok {
		s.cron.Remove(entryID)
		_, err := s.db.Exec(`UPDATE "public"."jobs" SET "cron" = $1 , "payload" = $2 WHERE "name" = $3 AND "cron" != '-'`, cron, payload, event)
		if err != nil {
			log.Print("schedule cron update error: ", err)
		}
	} else {
		_, err := s.db.Exec(`INSERT INTO "public"."jobs" ("name", "payload", "runAt", "cron") VALUES ($1, $2, $3, $4)`, event, payload, time.Now(), cron)
		if err != nil {
			log.Print("schedule cron insert error: ", err)
		}
	}

	eventFn, ok := s.listeners[event]
	if ok {
		entryID, err := s.cron.AddFunc(cron, func() { eventFn(payload) })
		s.cronEntries[event] = entryID
		if err != nil {
			log.Print("💀 error: ", err)
		}
	}
}

// attachCronJobs attaches cron jobs
func (s Scheduler) attachCronJobs() {
	log.Printf("Attaching cron jobs")
	rows, err := s.db.Query(`SELECT "id", "name", "payload", "cron" FROM "public"."jobs" WHERE "cron"!='-'`)
	if err != nil {
		log.Print("💀 error: ", err)
	}
	for rows.Next() {
		evt := Event{}
		rows.Scan(&evt.ID, &evt.Name, &evt.Payload, &evt.Cron)
		eventFn, ok := s.listeners[evt.Name]
		if ok {
			entryID, err := s.cron.AddFunc(evt.Cron, func() { eventFn(evt.Payload) })
			s.cronEntries[evt.Name] = entryID

			if err != nil {
				log.Print("💀 error: ", err)
			}
		}
	}
}

// StartCron starts cron job
func (s Scheduler) StartCron() func() {
	s.attachCronJobs()
	s.cron.Start()

	return func() {
		s.cron.Stop()
	}
}
