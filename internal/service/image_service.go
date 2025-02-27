package service

import (
	"path/filepath"
	"strings"

	"github.com/abhinandpn/CompressImage/server"
)

// ProcessAndCompressImage handles image processing via Imaginary API
func ProcessAndCompressImage(filename string, imageData []byte, size int64) (map[string]string, error) {
	// Generate a readable filename
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	baseName = strings.ReplaceAll(baseName, " ", "_")

	// Save original (apply reduction based on size)
	originalPath, err := server.ProcessImageWithImaginary(imageData, determineOriginalSizeReduction(size), baseName+"_original")
	if err != nil {
		return nil, err
	}

	// Save in different sizes
	sizes := map[string]int{
		"compressed_200-300KB": 250,
		"compressed_10-20KB":   15,
		"compressed_1-5KB":     3,
	}

	imagePaths := map[string]string{"original": originalPath}
	for key, quality := range sizes {
		path, err := server.ProcessImageWithImaginary(imageData, quality, baseName+"_"+key)
		if err == nil {
			imagePaths[key] = path
		}
	}

	return imagePaths, nil
}

// determineOriginalSizeReduction determines the quality based on file size
func determineOriginalSizeReduction(size int64) int {
	switch {
	case size > 5*1024*1024:
		return 50
	case size > 2*1024*1024:
		return 30
	default:
		return 100 // No compression
	}
}
