package stat

import (
	"context"
)

type Stat struct {
	Link string
	IP   string
}

type StatStore interface {
	Add(ctx context.Context, stat Stat) error
	ReadAll(ctx context.Context, shortLink string) (stats *[]Stat, err error)
	ReadIP(ctx context.Context, stat Stat) (count int64, err error)
}

type Stats struct {
	sstore StatStore
}

func NewStats(sstore StatStore) *Stats {
	return &Stats{
		sstore: sstore,
	}
}

func (s *Stats) Add(ctx context.Context, stat Stat) error {
	return s.sstore.Add(ctx, stat)
}

func (s *Stats) ReadAll(ctx context.Context, shortLink string) (stats *[]Stat, err error) {
	return s.sstore.ReadAll(ctx, shortLink)
}

func (s *Stats) ReadIP(ctx context.Context, stat Stat) (count int64, err error) {
	return s.sstore.ReadIP(ctx, stat)
}
