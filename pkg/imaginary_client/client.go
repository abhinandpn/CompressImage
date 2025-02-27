package imaginary_client

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/abhinandpn/CompressImage/internal/config"
	"github.com/abhinandpn/CompressImage/internal/repository"
)

// ResizeImage uses Imaginary API to compress images
func ResizeImage(file *multipart.FileHeader, quality int, outputName string) (string, error) {
	url := fmt.Sprintf("%s/resize?width=1000&height=1000&quality=%d", config.GetImaginaryURL(), quality)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	req, err := http.NewRequest("POST", url, src)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "image/jpeg")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Save file
	filePath, err := repository.SaveImageToStorage(outputName+".jpg", body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
