package usecase

import (
	"batikin-be/internal/app/motif/repository"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/domain/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

type MotifUsecaseItf interface {
	GetAll() ([]entity.Motif, error)
	GetSpecific(ctx *fiber.Ctx) (entity.Motif, error)
	Create(ctx *fiber.Ctx, req dto.CreateMotifRequest) (entity.Motif, error)
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
	motif := entity.Motif{
		Name:   req.Name,
		Prompt: req.Prompt,
	}

	prompt := `
		Generate a highly detailed, seamless, tileable batik pattern, designed for repeating across a surface without visible edges based on the prompt below.

		Prompt: ungu dengan bunga mawar yang mekar.
	`

	// Generate image
	imageResponse, err := u.openaiClient.Images.Generate(ctx.Context(), openai.ImageGenerateParams{
		Prompt:  prompt,
		N:       openai.Int(1),
		Model:   "gpt-image-1",
		Quality: "medium",
	})

	if err != nil {
		return entity.Motif{}, err
	}

	motif.ImageURL = imageResponse.Data[0].URL

	err = u.motifRepository.Create(&motif)
	if err != nil {
		return entity.Motif{}, err
	}

	response, err := u.motifRepository.GetSpecific(motif)
	if err != nil {
		return entity.Motif{}, err
	}

	return response, nil
}
