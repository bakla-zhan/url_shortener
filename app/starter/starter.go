package starter

import (
	"context"
	"sync"

	"github.com/bakla-zhan/url_shortener/app/repos/link"
	"github.com/bakla-zhan/url_shortener/app/repos/stat"
)

type App struct {
	Ls *link.Links
	Ss *stat.Stats
}

type AppStore interface {
	link.LinkStore
	stat.StatStore
}

func NewApp(as AppStore) *App {
	app := &App{
		Ls: link.NewLinks(as),
		Ss: stat.NewStats(as),
	}
	return app
}

type APIServer interface {
	Start(a *App)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a)
	<-ctx.Done()
	hs.Stop()
}
