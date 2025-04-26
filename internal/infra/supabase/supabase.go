package supabase

import (
	"batikin-be/config"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadImage(bucketName string, objectPath string, localFilePath string) (string, error) {
	conf := config.Load()
	supabaseURL := conf.SupabaseURL
	serviceRoleKey := conf.SupbaseToken

	fileBytes, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, objectPath)

	req, err := http.NewRequest("POST", url, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+serviceRoleKey)
	req.Header.Set("Content-Type", "image/png")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err

	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", string(respBody))

	return url, nil
}
