package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/JoYBoY1210/zephyr/utils"
)

func ServeFileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return

	}
	id := parts[2]
	sig := r.URL.Query().Get("sign")
	urlPath := fmt.Sprintf("/files/%s", id)
	if sig == "" || !utils.VerifyHMAC(urlPath, sig) {
		http.Error(w, "invalid or missing signature", http.StatusForbidden)
		return
	}
	var filepath, fileType string

	err := db.QueryRow("SELECT original_path,file_type FROM files WHERE id=?", id).Scan(&filepath, &fileType)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	f, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "error opening file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", fileType)
	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "error sending file", http.StatusInternalServerError)
		return
	}
}
