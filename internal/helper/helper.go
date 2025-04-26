package helper

import (
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"mime/multipart"
	"net/http"
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
