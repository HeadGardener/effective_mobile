package storage

import (
	"github.com/HeadGardener/effective_mobile/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB(conf config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", conf.URL)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
