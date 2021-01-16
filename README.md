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
2021/01/16 11:45:57 💾 Seeding database with table...
2021/01/16 11:45:57 🚀 Scheduling event SendEmail to run at 2021-01-16 11:46:57.177316418 +0545 +0545 m=+60.007979630
2021/01/16 11:45:57 🚀 Scheduling event PayBills to run at 2021-01-16 11:50:57.180342726 +0545 +0545 m=+300.011005907
2021/01/16 11:46:57 ⏰ Ticks Received...
2021/01/16 11:46:57 📨 Sending email with data:  mail: nilkantha.dipesh@gmail.com
2021/01/16 11:47:57 ⏰ Ticks Received...
^C2021/01/16 11:48:00 
❌ Interrupt received closing...
```