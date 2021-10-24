package starter

import (
	"context"
	"sync"
	"urlshortener/app/repo/link"
)

type App struct {
	Ls *link.Links
}

func NewApp(l link.LinkStore) *App {
	a := &App{
		Ls: link.NewLinks(l),
	}
	return a
}

type APIServer interface {
	Start(ls *link.Links)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a.Ls)
	<-ctx.Done()
	hs.Stop()
}
