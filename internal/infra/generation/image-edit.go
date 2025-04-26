package generation

import (
	"batikin-be/config"
	"batikin-be/internal/infra/supabase"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type Response struct {
	Data []struct {
		Base64 string `json:"b64_json"`
	}
}

func GenerateImageEdit(imagePath string, prompt string) (string, error) {
	client := resty.New()
	conf := config.Load()
	url := "https://api.openai.com/v1/images/edits"

	resp, err := client.R().
		SetFile("image", imagePath).
		SetFormData(map[string]string{
			"prompt": prompt,
			"model":  "gpt-image-1",
			"size":   "1024x1024",
		}).
		SetHeader("Authorization", "Bearer "+conf.OPENAIAPIKey).
		Post(url)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		fmt.Println("Response Status:", resp)
		return "", errors.New("failed to parse file")
	}

	var response Response
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return "", err
	}

	imageData, err := base64.StdEncoding.DecodeString(response.Data[0].Base64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}

	uniqueFileName := uuid.New().String() + ".png"
	path := "./tmp/" + uniqueFileName

	err = os.WriteFile(path, imageData, os.ModePerm)
	if err != nil {
		return "", err
	}

	publicUrl, err := supabase.UploadImage("motif", uniqueFileName, path)
	if err != nil {
		return "", err
	}

	err = os.Remove(path)
	if err != nil {
		return "", err
	}

	err = os.Remove(imagePath)
	if err != nil {
		return "", err
	}

	return string(publicUrl), nil
}
