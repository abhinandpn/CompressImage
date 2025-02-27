package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/abhinandpn/CompressImage/internal/service"
)

// UploadImageHandler handles image uploads
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the uploaded file
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file size limit (10MB)
	if header.Size > 10*1024*1024 {
		http.Error(w, "File size exceeds 10MB", http.StatusBadRequest)
		return
	}

	// Read the file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Process image
	imagePaths, err := service.ProcessAndCompressImage(header.Filename, fileBytes, header.Size)
	if err != nil {
		http.Error(w, "Failed to process image", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image uploaded successfully",
		"paths":   imagePaths,
	})
}
