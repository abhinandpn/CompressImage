package service

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/abhinandpn/CompressImage/server"
)

// ProcessAndCompressImage handles image processing via Imaginary API (concurrent)
func ProcessAndCompressImage(filename string, imageData []byte, size int64) (map[string]string, error) {
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	baseName = strings.ReplaceAll(baseName, " ", "_")

	if cachedPaths, exists := GetCachedResult(baseName); exists {
		return cachedPaths, nil
	}

	var wg sync.WaitGroup
	resultChan := make(chan struct {
		key  string
		path string
		err  error
	}, 4)

	sizes := map[string]int{
		"original":             determineOriginalSizeReduction(size), // ✅ Fixed: Now defined
		"compressed_200-300KB": 250,
		"compressed_10-20KB":   15,
		"compressed_1-5KB":     3,
	}

	for key, quality := range sizes {
		wg.Add(1)
		go func(k string, q int) {
			defer wg.Done()
			path, err := server.ProcessImageWithImaginary(imageData, q, baseName+"_"+k)
			resultChan <- struct {
				key  string
				path string
				err  error
			}{key: k, path: path, err: err}
		}(key, quality)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	imagePaths := make(map[string]string)
	for res := range resultChan {
		if res.err == nil {
			imagePaths[res.key] = res.path
		}
	}

	CacheResult(baseName, imagePaths)
	return imagePaths, nil
}

// ✅ Added function to determine compression based on file size
func determineOriginalSizeReduction(size int64) int {
	switch {
	case size > 5*1024*1024: // If size > 5MB
		return 50
	case size > 2*1024*1024: // If size is between 2MB and 5MB
		return 30
	default: // If size is < 2MB
		return 100 // No compression
	}
}
