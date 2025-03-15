package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abhinandpn/CompressImage/internal/config"
	handler "github.com/abhinandpn/CompressImage/internal/handler" // ✅ Import the handler package
	"github.com/abhinandpn/CompressImage/server"
)

func main() {
	// Start Imaginary server in the background
	server.StartImaginaryServer()
	// Load environment variables
	config.LoadEnv()
	// Create storage directory if it doesn't exist
	err := os.MkdirAll("storage", os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create storage directory:", err)
	}

	// Register HTTP handlers
	http.HandleFunc("/upload", handler.UploadImageHandler) // ✅ Now handler is recognized
	http.HandleFunc("/s3upload", handler.S3ImageHandler)   // ✅ Now handler is recognized

	port := "3000"
	fmt.Println("Server running on port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
