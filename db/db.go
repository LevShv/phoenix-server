package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal("Ошибка открытия БД:", err)
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatal("Ошибка поиска файлов:", err)
	}

	for _, file := range files {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Ошибка чтения файла %s: %v", file, err)
		}

		_, err = DB.Exec(string(sqlBytes))
		if err != nil {
			log.Fatalf("Ошибка выполнения sql %s: %v", file, err)
		}

		log.Printf("✓ Выполнен: %s", file)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
