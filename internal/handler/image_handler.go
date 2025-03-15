package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg" // Import for JPEG decoding
	_ "image/png"  // Import for PNG decoding
	"io"
	"io/ioutil"
	"net/http"

	"github.com/abhinandpn/CompressImage/internal/service"
)

// UploadImageHandler handles multiple image uploads
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit per file
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["image"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var imagesData []map[string]interface{}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if fileHeader.Size > 10*1024*1024 {
			http.Error(w, "File size exceeds 10MB", http.StatusBadRequest)
			return
		}

		// Read file into memory
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Decode image to get dimensions
		imgConfig, _, err := image.DecodeConfig(bytes.NewReader(fileBytes))
		if err != nil {
			http.Error(w, "Failed to decode image", http.StatusInternalServerError)
			return
		}
		originalWidth := imgConfig.Width
		originalHeight := imgConfig.Height

		// Process and compress image with aspect ratio preservation
		imagePaths, err := service.ProcessAndCompressImage(fileHeader.Filename, fileBytes, fileHeader.Size, originalWidth, originalHeight)
		if err != nil {
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}

		// Calculate aspect ratio without max=10 constraint
		aspectRatio := calculateAspectRatio(originalWidth, originalHeight)

		// Append image data
		imagesData = append(imagesData, map[string]interface{}{
			"filename":        fileHeader.Filename,
			"aspect_ratio":    aspectRatio,
			"original_width":  originalWidth,
			"original_height": originalHeight,
			"paths":           imagePaths,
		})
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Images uploaded successfully",
		"images":  imagesData,
	})
}

// Function to calculate the greatest common divisor (GCD)
func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// Function to calculate and simplify aspect ratio (without max=10 limitation)
func calculateAspectRatio(width, height int) string {
	if width == 0 || height == 0 {
		return "Invalid"
	}

	g := gcd(width, height) // Find the GCD
	simplifiedWidth := width / g
	simplifiedHeight := height / g

	return fmt.Sprintf("%d:%d", simplifiedWidth, simplifiedHeight)
}

// S3ImageHandler handles image upload, processing, and saving to S3
func S3ImageHandler(w http.ResponseWriter, r *http.Request) {
	// Limit the file size to prevent too large uploads
	const maxFileSize = 10 * 1024 * 1024 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	// Parse the form and get the file
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file into a byte slice
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Get the original image dimensions (width and height)
	originalImage, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return
	}
	originalWidth := originalImage.Bounds().Dx()
	originalHeight := originalImage.Bounds().Dy()

	// Call S3ProcessAndCompressImage to process and upload the image to S3
	imagePaths, err := service.S3ProcessAndCompressImage(fileHeader.Filename, fileBytes, fileHeader.Size, originalWidth, originalHeight)
	if err != nil {
		http.Error(w, "Failed to process and upload image to S3", http.StatusInternalServerError)
		return
	}

	// Check if any images were uploaded successfully
	if len(imagePaths) == 0 {
		http.Error(w, "No images were uploaded", http.StatusInternalServerError)
		return
	}

	// Return a success response with the S3 URLs of the uploaded images
	response := map[string]interface{}{
		"message":   "Image processed and uploaded successfully",
		"imageUrls": imagePaths,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
