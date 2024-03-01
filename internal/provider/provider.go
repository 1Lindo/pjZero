package provider

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPgConn(storagePath string) (*sqlx.DB, error) {
	const op string = "storage.provider.postgres.newConn" //op == operation

	pg, err := sqlx.Open("postgres", storagePath)
	if err != nil {
		fmt.Printf("%s:%w", op, err)
		return nil, err
	}

	stmtTable, err := pg.Prepare(`
	CREATE TABLE IF NOT EXISTS videos(
    	id SERIAL PRIMARY KEY,
    	file_name TEXT NOT NULL,
	    file_path TEXT NOT NULL
	    );
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmtTable.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return pg, nil
}
