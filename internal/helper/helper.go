package helper

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandomNumber(n int) string {
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func StringISOToDateTime(dateString string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateString)
}

func GenerateSignature(orderID, statusCode, grossAmount, serverKey string) string {
	signature := orderID + statusCode + grossAmount + serverKey
	hash := sha512.New()
	hash.Write([]byte(signature))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetFileMimeType(file *multipart.FileHeader) (string, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Read a small portion of the file (first 512 bytes)
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}

	// Detect the MIME type
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}

func getFilenameFromURL(url string) string {
	// Split the URL by "/" and get the last part
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		filename := parts[len(parts)-1]

		// Remove any URL parameters
		if queryIndex := strings.Index(filename, "?"); queryIndex != -1 {
			filename = filename[:queryIndex]
		}

		// If the filename is not empty, return it
		if filename != "" {
			return filename
		}
	}
	return ""
}

func FetchAndSaveImage(url, folderPath string) (io.Reader, string, error) {
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, "", fmt.Errorf("failed to create directory: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("non-success status code: %d", resp.StatusCode)
	}

	filename := getFilenameFromURL(url)
	if filename == "" {
		filename = "image.jpg" // Default filename if none could be determined
	}

	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)
	uniqueFilename := fmt.Sprintf("%s_%d%s", basename, timestamp, ext)

	filePath := filepath.Join(folderPath, uniqueFilename)

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	err = os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		// Try alternate location if permission denied
		if os.IsPermission(err) {
			altPath := filepath.Join(os.TempDir(), uniqueFilename)
			err = os.WriteFile(altPath, imageData, 0644)
			if err != nil {
				return nil, "", fmt.Errorf("failed to write file to alternate location: %w", err)
			}
			filePath = altPath
		} else {
			return nil, "", fmt.Errorf("failed to write file: %w", err)
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open saved file: %w", err)
	}

	return file, filePath, nil
}
