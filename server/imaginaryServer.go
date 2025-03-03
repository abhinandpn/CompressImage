package server

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/exec"

	"github.com/nfnt/resize"
)

// ProcessImageWithImaginary calls Imaginary API to resize/compress images
// ProcessImageWithImaginary compresses and resizes an image while keeping aspect ratio
func ProcessImageWithImaginary(imageData []byte, quality int, outputName string, width int, height int) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", err
	}

	// Resize the image while keeping the aspect ratio
	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	// Create the output file
	outputPath := fmt.Sprintf("storage/%s.jpg", outputName)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Encode the resized image
	options := &jpeg.Options{Quality: quality}
	err = jpeg.Encode(outFile, resizedImg, options)
	if err != nil {
		return "", err
	}

	return outputPath, nil
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
