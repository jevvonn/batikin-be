package usecase

import (
	"batikin-be/internal/app/motif/repository"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/domain/entity"
	"batikin-be/internal/infra/generation"
	"batikin-be/internal/infra/supabase"
	"fmt"
	"io"
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

	return dto.CaptureMotifResponse{
		ImageURL: url,
	}, nil
}
