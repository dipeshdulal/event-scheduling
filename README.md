### Database Based Event Scheduling

Example that demonstrates super basic database based event scheduling.


#### To run this example; 
- Copy `.env.example` to `.env` and update postgres database dsn.
- Run `go run .` 


**Basically,**
 - When we want to schedule a job, we add to database 
```go
scheduler.Schedule("SendEmail", "mail: nilkantha.dipesh@gmail.com", time.Now().Add(1*time.Minute)) 
```

- Another go routine is always looking for jobs to execute (that has time expired) in the given interval.
```go
scheduler.CheckEventsInInterval(ctx, time.Minute)
```

**Output looks like;**
```
2021/01/16 11:58:49 💾 Seeding database with table...
2021/01/16 11:58:49 🚀 Scheduling event SendEmail to run at 2021-01-16 11:59:49.344904505 +0545 +0545 m=+60.004623549
2021/01/16 11:58:49 🚀 Scheduling event PayBills to run at 2021-01-16 12:00:49.34773798 +0545 +0545 m=+120.007457039
2021/01/16 11:59:49 ⏰ Ticks Received...
2021/01/16 11:59:49 📨 Sending email with data:  mail: nilkantha.dipesh@gmail.com
2021/01/16 12:00:49 ⏰ Ticks Received...
2021/01/16 12:01:49 ⏰ Ticks Received...
2021/01/16 12:01:49 💲 Pay me a bill:  paybills: $4,000 bill
2021/01/16 12:02:49 ⏰ Ticks Received...
2021/01/16 12:03:49 ⏰ Ticks Received...
^C2021/01/16 12:03:57 
❌ Interrupt received closing...
```