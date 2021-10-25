package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"urlshortener/api/handlers"
	"urlshortener/api/server"
	"urlshortener/app/starter"
	"urlshortener/db/sql/pg"
)

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	// output current time zone
	tnow := time.Now()
	tz, _ := tnow.Zone()

	dsn, ok := os.LookupEnv("DB_DSN")
	if !ok {
		log.Fatal("DB_DSN env is not set")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	ls, err := pg.NewLinks(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer ls.Close()

	a := starter.NewApp(ls)
	h := handlers.NewHandlers(a.Ls)
	srv := server.NewServer(":8080", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	log.Printf("Local time zone %s. Service started at %s", tz, tnow.Format("2006-01-02T15:04:05.000 MST"))

	<-ctx.Done()
	cancel()
	wg.Wait()
}
