package utils

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	// Cek koneksi ke database
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
