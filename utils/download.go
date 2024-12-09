package utils

import (
	"fmt"
	"image"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// generateRandomString generate random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		result.WriteByte(charset[r.Intn(len(charset))])
	}
	return result.String()
}

func getFileExtensionFromMIME(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "application/pdf":
		return ".pdf"
	case "text/plain":
		return ".txt"
	default:
		return "" // If no known extension is found, return empty
	}
}

// DownloadFile downloads URL file and save
func DownloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Get the filename from the Content-Disposition header (if available)
	var filename string
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		parts := strings.Split(contentDisposition, "filename=")
		if len(parts) > 1 {
			filename = strings.Trim(parts[1], `"`)
		}
	}

	// If filename is not set generate random filename from string charset
	if filename == "" {
		filename = generateRandomString(16) // generate a random name
	}

	if !strings.Contains(filename, ".") {
		contentType := resp.Header.Get("Content-Type")
		extension := getFileExtensionFromMIME(contentType)
		if extension != "" {
			filename = filename + extension
		} else {
			filename = filename + ".bin"
		}
	}

	// Ensure the ./downloads directory exists
	dir := "./downloads"
	err = os.MkdirAll(dir, os.ModePerm) // Create the directory if it doesn't exist
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Full file path in the ./downloads directory
	filePath := filepath.Join(dir, filename)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write content to file: %v", err)
	}

	return filePath, nil
}

type ImageDetailsType struct {
	Height int
	Width  int
	Size   int64
}

// ImageDetails calculate image width and height
func ImageDetails(src string) (*ImageDetailsType, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	m, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %v", err)
	}
	img := &ImageDetailsType{
		Height: m.Bounds().Dy(),
		Width:  m.Bounds().Dx(),
		Size:   stat.Size(),
	}

	return img, nil
}
