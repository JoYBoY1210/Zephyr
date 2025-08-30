package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/JoYBoY1210/zephyr/utils"
)

type signupReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var signup signupReq
	if err := json.NewDecoder(r.Body).Decode(&signup); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if signup.Email == "" || signup.Password == "" || signup.Username == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	hash, err := utils.HashPassword(signup.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}
	res, err := db.Exec(`INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`, signup.Username, signup.Email, hash)
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}
	userId, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "error retrieving user id", http.StatusInternalServerError)
		return
	}
	token, err := utils.GenerateToken(32)
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}
	expires := time.Now().Add(7 * 24 * time.Hour)
	_, err = db.Exec(
		`INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)`,
		userId, token, expires,
	)
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "zephyr_session",
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user created and logged in"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var login loginReq
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if login.Email == "" || login.Password == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	var userId int64
	var hash string
	err := db.QueryRow(`SELECT id, password_hash FROM users WHERE email=?`, login.Email).Scan(&userId, &hash)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	if !utils.CheckPassword(login.Password, hash) {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateToken(32)
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}
	expires := time.Now().Add(7 * 24 * time.Hour)
	_, err = db.Exec(
		`INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)`,
		userId, token, expires,
	)
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "zephyr_session",
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	w.Write([]byte("login successful"))
}
