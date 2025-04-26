package usecase

import (
	"batikin-be/config"
	"batikin-be/internal/app/motif/repository"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/domain/entity"
	"batikin-be/internal/infra/generation"
	"batikin-be/internal/infra/supabase"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

type MotifUsecaseItf interface {
	GetAll() ([]entity.Motif, error)
	GetSpecific(ctx *fiber.Ctx) (entity.Motif, error)
	Create(ctx *fiber.Ctx, req dto.CreateMotifRequest) (entity.Motif, error)
	Capture(ctx *fiber.Ctx) (dto.CaptureMotifResponse, error)
}

type MotifUsecase struct {
	motifRepository repository.MotifPostgreSQLItf
	openaiClient    openai.Client
}

func NewMotifUsecase(motifRepository repository.MotifPostgreSQLItf, openaiClient openai.Client) MotifUsecaseItf {
	return &MotifUsecase{motifRepository, openaiClient}
}

func (u *MotifUsecase) GetAll() ([]entity.Motif, error) {
	return u.motifRepository.GetAll()
}

func (u *MotifUsecase) GetSpecific(ctx *fiber.Ctx) (entity.Motif, error) {
	param := ctx.Params("id")
	motifId, err := uuid.Parse(param)
	if err != nil {
		return entity.Motif{}, err
	}

	return u.motifRepository.GetSpecific(entity.Motif{
		ID: motifId,
	})
}

func (u *MotifUsecase) Create(ctx *fiber.Ctx, req dto.CreateMotifRequest) (entity.Motif, error) {
	userId := ctx.Locals("userId").(string)
	motif := &entity.Motif{
		Prompt: req.Prompt,
		UserID: uuid.MustParse(userId),
	}

	prompt := fmt.Sprintf(`generate an image batik pattern with size of 1024x1024 based on this description :

		%s

		The image should be a seamless batik pattern
	`, req.Prompt)

	// Generate image
	url, err := generation.GenerateImage(prompt)
	if err != nil {
		return entity.Motif{}, err
	}
	motif.ImageURL = url

	if req.Name == "" {
		prompt = fmt.Sprintf(`generate a single title without any description for this batik pattern based on this description : %s`, req.Prompt)

		chatCompletion, err := u.openaiClient.Chat.Completions.New(ctx.Context(), openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model: openai.ChatModelGPT4o,
		})
		if err != nil {
			return entity.Motif{}, err
		}

		motif.Name = strings.ReplaceAll(chatCompletion.Choices[0].Message.Content, "\"", "")
	} else {
		motif.Name = req.Name
	}

	err = u.motifRepository.Create(motif)
	if err != nil {
		return entity.Motif{}, err
	}

	response, err := u.motifRepository.GetSpecific(*motif)
	if err != nil {
		return entity.Motif{}, err
	}

	return response, nil
}

func (u *MotifUsecase) Capture(ctx *fiber.Ctx) (dto.CaptureMotifResponse, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("file not found")
	}

	uniqueFileName := uuid.New().String() + path.Ext(file.Filename)
	path := "./tmp/" + uniqueFileName

	fileData, err := file.Open()
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to open file")
	}
	defer fileData.Close()

	data, err := io.ReadAll(fileData)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to read file")
	}

	err = os.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to write file")
	}

	url, err := supabase.UploadImage("capture", uniqueFileName, path)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to upload file")
	}

	err = os.Remove(path)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to remove file")
	}

	// Get roboflow
	conf := config.Load()
	reqBody := map[string]interface{}{
		"api_key": conf.RoboflowAPIKey,
		"inputs": map[string]interface{}{
			"image": map[string]interface{}{
				"type":  "url",
				"value": url,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to encode json")
	}

	apiURL := "https://serverless.roboflow.com/infer/workflows/batikkin/detect-and-classify"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to send request")
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to read response")
	}

	var result struct {
		Outputs []struct {
			OutputImage struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"output_image"`
		} `json:"outputs"`
	}

	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to unmarshal response")
	}

	if len(result.Outputs) == 0 {
		return dto.CaptureMotifResponse{}, fmt.Errorf("no outputs found")
	}

	imageData, err := base64.StdEncoding.DecodeString(result.Outputs[0].OutputImage.Value)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to decode base64 image")
	}

	uniqueFileName = uuid.New().String() + ".png"
	path = "./tmp/" + uniqueFileName

	err = os.WriteFile(path, imageData, os.ModePerm)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to write file")
	}

	publicUrl, err := supabase.UploadImage("scoring", uniqueFileName, path)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to upload file")
	}

	err = os.Remove(path)
	if err != nil {
		return dto.CaptureMotifResponse{}, fmt.Errorf("failed to remove file")
	}

	return dto.CaptureMotifResponse{
		ImageURL: publicUrl,
	}, nil
}
