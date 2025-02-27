package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"

	"github.com/abhinandpn/CompressImage/internal/repository"
)

// ProcessImageWithImaginary calls Imaginary API to resize/compress images
func ProcessImageWithImaginary(imageData []byte, quality int, outputName string) (string, error) {
	url := fmt.Sprintf("http://localhost:9000/resize?width=1000&height=1000&quality=%d", quality)

	// Create a new HTTP request with multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", outputName+".jpg")
	if err != nil {
		return "", err
	}
	part.Write(imageData)
	writer.Close()

	// Send request to Imaginary
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := HttpClient // Use the optimized HTTP client
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response (processed image)
	resizedImage, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Save the resized image to storage
	filePath, err := repository.SaveImageToStorage(outputName+".jpg", resizedImage)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// StartImaginaryServer starts the Imaginary server as a background process
func StartImaginaryServer() {
	// Ensure `imaginary` is installed
	_, err := exec.LookPath("imaginary")
	if err != nil {
		log.Fatal("Imaginary is not installed. Install it using: go install github.com/h2non/imaginary@latest")
	}

	// Run Imaginary server in the background
	cmd := exec.Command("imaginary", "-p", "9000", "-enable-url-source")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		fmt.Println("Starting Imaginary server on port 9000...")
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to start Imaginary server: %v", err)
		}
	}()
}
