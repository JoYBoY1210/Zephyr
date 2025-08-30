package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() *sql.DB {
	var err error
	DB, err = sql.Open("sqlite3", "zephyr.db")
	if err != nil {
		log.Fatal("failed to connect to the database: ", err)
	}

	fileTable := `
	CREATE TABLE IF NOT EXISTS files(
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  uuid TEXT NOT NULL UNIQUE
	  filename TEXT NOT NULL,
	  file_type TEXT NOT NULL,
	  upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
	  original_path TEXT NOT NULL
	);
	`


	_, err = DB.Exec(fileTable)
	if err != nil {
		log.Fatal("failed to create files table: ", err)
	}
	
	log.Println("database initialized and files table done")
	return DB
}
