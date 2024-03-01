package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

//TODO: перенести описание интерфейса на место его использования

type IPgRepo interface {
	InsertData(fileName string, filePath string) error
	GetData(filePath string) (string, error)
}

var (
	ErrUrlNotound = "Url is not found"
	ErrUrlExists  = "Url already exists"
)

type pgRepo struct {
	db *sqlx.DB
}

func (p *pgRepo) InsertData(fileName string, filePath string) error {
	const op = "internal.repository.postgres.insertVideoData"

	stmt, err := p.db.Prepare(`INSERT INTO videos (file_name, file_path) VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(fileName, filePath)
	if err != nil {
		//TODO: доработать на случай если захотим добавить что-то с уже исп alias
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	// в MySQL такого нет поэтому стоит убрать int64 в return
	//id, err := res.LastInsertId()
	//if err != nil {
	//	return fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	//}
	//fmt.Println(id)

	return nil
}

func (p *pgRepo) GetData(fileName string) (string, error) {
	const op = "internal.repository.postgres.getVideoData"

	row := p.db.QueryRow(`SELECT file_path FROM videos WHERE file_name = $1`, fileName)

	var resFilePath string
	err := row.Scan(&resFilePath)
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resFilePath, nil
}
