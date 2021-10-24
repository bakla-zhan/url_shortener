package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"urlshortener/api/handlers"
	"urlshortener/api/server"
	"urlshortener/app/starter"
	"urlshortener/db/sql/pg"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	ls, err := pg.NewLinks("postgres://postgres:1@192.168.1.138/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer ls.Close()

	// lst := link.NewLinks(ls)
	a := starter.NewApp(ls)
	h := handlers.NewHandlers(a.Ls)
	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
