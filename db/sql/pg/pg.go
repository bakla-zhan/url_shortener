package pg

import (
	"context"
	"database/sql"
	"log"
	"urlshortener/app/repo/link"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var _ link.LinkStore = &Links{}

type Links struct {
	db *sql.DB
}

func NewLinks(dsn string) (*Links, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	l := &Links{
		db: db,
	}
	return l, nil
}

func (l *Links) Close() {
	l.db.Close()
}

func (l *Links) Create(ctx context.Context, link link.Link) error {
	_, err := l.db.ExecContext(ctx, `INSERT INTO links (short, long) values ($1, $2)`, link.Short, link.Long)
	if err != nil {
		return err
	}
	return nil
}

func (l *Links) Read(ctx context.Context, shortLink string) (longLink string, err error) {
	if err = l.db.QueryRowContext(ctx, `SELECT long FROM links WHERE short = $1`, shortLink).Scan(&longLink); err != nil {
		log.Println("db read", err)
		return "", err
	}
	return longLink, nil
}
