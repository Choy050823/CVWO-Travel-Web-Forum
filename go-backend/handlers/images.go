package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const uploadDir = "./uploads"

func init() {
	// Create upload directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
}

func UploadImages(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parse multipart form (max 10MB files)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get files from request
	files := r.MultipartForm.File["images"]
	var imageUrls []string

	for _, fileHeader := range files {
		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create unique filename
		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy file content
		if _, err := io.Copy(dst, file); err != nil  {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		imageUrls = append(imageUrls, "/"+filePath)
	}

	// Return image URLs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"imageUrls": imageUrls,
	})
}