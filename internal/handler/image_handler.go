package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/abhinandpn/CompressImage/internal/service"
)

// UploadImageHandler handles image uploads
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 10*1024*1024 {
		http.Error(w, "File size exceeds 10MB", http.StatusBadRequest)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	imagePaths, err := service.ProcessAndCompressImage(header.Filename, fileBytes, header.Size)
	if err != nil {
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image uploaded successfully",
		"paths":   imagePaths,
	})
}
