package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handler "github.com/abhinandpn/CompressImage/internal/handler" // ✅ Import the handler package
	"github.com/abhinandpn/CompressImage/server"
)

func main() {
	// Start Imaginary server in the background
	server.StartImaginaryServer()

	// Create storage directory if it doesn't exist
	err := os.MkdirAll("storage", os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create storage directory:", err)
	}

	// Register HTTP handlers
	http.HandleFunc("/upload", handler.UploadImageHandler) // ✅ Now handler is recognized

	port := "8080"
	fmt.Println("Server running on port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
