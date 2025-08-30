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
	  user_id INTEGER NOT NULL,
	  filename TEXT NOT NULL,
	  file_type TEXT NOT NULL,
	  upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
	  original_path TEXT NOT NULL
	);
	`

	userTable := `
	CREATE TABLE IF NOT EXISTS users(
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  username TEXT NOT NULL,
	  email TEXT UNIQUE NOT NULL,
	  password_hash TEXT NOT NULL,
	  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	sessionTable := `
	CREATE TABLE IF NOT EXISTS sessions(
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  user_id INTEGER NOT NULL,
	  token TEXT NOT NULL,
	  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	  expires_at DATETIME NOT NULL
	  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	  );`

	_, err = DB.Exec(fileTable)
	if err != nil {
		log.Fatal("failed to create files table: ", err)
	}
	_, err = DB.Exec(userTable)
	if err != nil {
		log.Fatal("failed to create user table: ", err)
	}
	_, err = DB.Exec(sessionTable)
	if err != nil {
		log.Fatal("failed to create session table: ", err)
	}
	log.Println("database initialized and files table done")
	return DB
}
