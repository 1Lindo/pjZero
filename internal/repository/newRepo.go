package repository

import "github.com/jmoiron/sqlx"

func NewPgRepo(pg *sqlx.DB) IPgRepo {
	return &pgRepo{
		db: pg,
	}
}
