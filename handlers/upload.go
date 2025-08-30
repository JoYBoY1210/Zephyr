package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/JoYBoY1210/zephyr/utils"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll("uploads", os.ModePerm)

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := filepath.Join("uploads", filename)

	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "couldnot save the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	filetype := detectFiletype(header)

	res, err := db.Exec(
		"INSERT INTO files (user_id, filename, file_type, original_path) VALUES (?, ?, ?, ?)",
		1,
		header.Filename,
		filetype,
		filepath,
	)
	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	fileId, _ := res.LastInsertId()

	urlPath := fmt.Sprintf("/files/%d", fileId)
	signature := utils.GenerateHMAC(urlPath)
	signedURL := fmt.Sprintf("%s?sign=%s", urlPath, signature)
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(signedURL))

}

func detectFiletype(header *multipart.FileHeader) string {
	if header.Header.Get("Content-Type") != "" {
		return header.Header.Get("Content-Type")
	}

	return "application/octet-stream"
}
