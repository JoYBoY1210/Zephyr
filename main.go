package main

import (
	"log"
	"net/http"

	"github.com/JoYBoY1210/zephyr/db"
	"github.com/JoYBoY1210/zephyr/handlers"
)

func main() {
	database := db.InitDB()
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadHandler(w, r, database)
	})
	http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeFileHandler(w, r, database)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, database)
	})
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handlers.SignupHandler(w, r, database)
	})
	log.Println("Server running at port 8080")
	http.ListenAndServe(":8080", nil)

}
