package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bakla-zhan/url_shortener/api/handlers"
	"github.com/bakla-zhan/url_shortener/api/server"
	"github.com/bakla-zhan/url_shortener/app/starter"
	"github.com/bakla-zhan/url_shortener/db/sql/pg"
)

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("DATABASE_URL env is not set")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT env is not set")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	st, err := pg.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	a := starter.NewApp(st)
	h := handlers.NewHandlers(a)
	srv := server.NewServer(fmt.Sprint(":", port), h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	tnow := time.Now()
	tz, _ := tnow.Zone()
	log.Printf("Local time zone %s. Service started at %s", tz, tnow.Format("2006-01-02T15:04:05.000 MST"))

	<-ctx.Done()
	cancel()
	wg.Wait()
}
