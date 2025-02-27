package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abhinandpn/CompressImage/internal/handler"
	"github.com/abhinandpn/CompressImage/server"
)

func main() {
	// Start Imaginary server in background
	server.StartImaginaryServer()

	// Create storage directory if not exists
	err := os.MkdirAll("storage", os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create storage directory:", err)
	}

	// Register HTTP handlers
	http.HandleFunc("/upload", handler.UploadImageHandler)

	port := "8080"
	fmt.Println("Server running on port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
