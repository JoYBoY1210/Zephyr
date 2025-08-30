package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB,storagePath,host string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(50 << 20)
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

	filetype:=detectFileType(header.Filename,header.Header.Get("Content-Type"))
	dir:=filepath.Join(storagePath,filetype)
	os.Mkdir(dir,os.ModePerm)
	u:=uuid.New().String()
	ext:=filepath.Ext(header.Filename)
	filename:=fmt.Sprintf("%s%s",u,ext)
	fullpath:=filepath.Join(dir,filename)

	

	dst, err := os.Create(fullpath)
	if err != nil {
		http.Error(w, "couldnot save the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	_, err = db.Exec(
			"INSERT INTO files (uuid, filename, file_type, original_path) VALUES (?, ?, ?, ?)",
			u, header.Filename, filetype, fmt.Sprintf("%s/%s", filetype, filename),
	)
	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}


	urlPath := fmt.Sprintf("%s/f/%s", strings.TrimRight(host, "/"), u)
	
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(urlPath))

}

func detectFileType(filename, contentType string) string {
	lc := strings.ToLower(filename)
	if strings.HasPrefix(contentType, "image/") || strings.HasSuffix(lc, ".png") || strings.HasSuffix(lc, ".jpg") || strings.HasSuffix(lc, ".jpeg") || strings.HasSuffix(lc, ".gif") {
		return "image"
	} else if strings.HasPrefix(contentType, "video/") || strings.HasSuffix(lc, ".mp4") || strings.HasSuffix(lc, ".mov") || strings.HasSuffix(lc, ".avi") {
		return "video"
	}
	return "other"
}
