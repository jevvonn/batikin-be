package generation

import (
	"batikin-be/config"
	"batikin-be/internal/infra/supabase"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func GenerateImage(prompt string) (string, error) {
	conf := config.Load()
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/imagen-3.0-generate-002:predict?key=%s", conf.GeminiAPIKey)

	requestBody := map[string]interface{}{
		"instances": []map[string]string{
			{
				"prompt": prompt,
			},
		},
		"parameters": map[string]interface{}{
			"sampleCount": 1,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Kirim request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Baca response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response Status:", resp.Status)

	// Parsing JSON response
	var jsonResponse struct {
		Predictions []struct {
			MimeType           string `json:"mimeType"`
			BytesBase64Encoded string `json:"bytesBase64Encoded"`
		} `json:"predictions"`
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return "", err
	}

	if len(jsonResponse.Predictions) == 0 {
		return "", fmt.Errorf("no predictions found in the response")
	}

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(jsonResponse.Predictions[0].BytesBase64Encoded)
	if err != nil {
		return "", err
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

	return publicUrl, nil
}
