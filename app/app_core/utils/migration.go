package utils

import (
	"database/sql"
	"log"
	"os"
)

func Migrate(db *sql.DB, migrationFile string) error {
	// Membaca file migrasi
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return err
	}

	// Menjalankan query migrasi
	_, err = db.Exec(string(content))
	if err != nil {
		return err
	}

	log.Println("Migration executed successfully")
	return nil
}
