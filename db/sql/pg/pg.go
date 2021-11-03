package pg

import (
	"context"
	"database/sql"

	"github.com/bakla-zhan/url_shortener/app/repos/link"
	"github.com/bakla-zhan/url_shortener/app/repos/stat"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var _ link.LinkStore = &Store{}
var _ stat.StatStore = &Store{}

type Store struct {
	db *sql.DB
}

func NewStore(dsn string) (*Store, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`create table IF NOT EXISTS links (short varchar(8), long varchar(1000));
		create index IF NOT EXISTS concurrently links_short_idx on links using btree (short text_pattern_ops);
	
		create table IF NOT EXISTS stats (link varchar(8), ip inet);
		create index IF NOT EXISTS concurrently stats_link_ip_idx on stats using btree (link text_pattern_ops, ip inet_ops);`)
	if err != nil {
		db.Close()
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	s := &Store{
		db: db,
	}
	return s, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) CreateLink(ctx context.Context, link link.Link) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO links (short, long) values ($1, $2)`, link.Short, link.Long)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ReadLink(ctx context.Context, shortLink string) (longLink string, err error) {
	if err = s.db.QueryRowContext(ctx, `SELECT long FROM links WHERE short = $1`, shortLink).Scan(&longLink); err != nil {
		return "", err
	}
	return longLink, nil
}

func (s *Store) Add(ctx context.Context, stat stat.Stat) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO stats (link, ip) values ($1, $2)`, stat.Link, stat.IP)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) ReadAll(ctx context.Context, shortLink string) (stats *[]stat.Stat, err error) {
	st := stat.Stat{}
	result := make([]stat.Stat, 0)

	rows, err := s.db.QueryContext(ctx, `SELECT link, ip FROM stats WHERE link = $1`, shortLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&st.Link,
			&st.IP,
		); err != nil {
			return nil, err
		}
		result = append(result, st)
	}

	return &result, nil
}
func (s *Store) ReadIP(ctx context.Context, stat stat.Stat) (count int64, err error) {
	result, err := s.db.ExecContext(ctx, `SELECT * FROM stats WHERE link = $1 and ip = $2`, stat.Link, stat.IP)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
