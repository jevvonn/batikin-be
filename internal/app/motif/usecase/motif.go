package usecase

import (
	"batikin-be/internal/app/motif/repository"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/domain/entity"
	"batikin-be/internal/infra/generation"
	"fmt"

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
	userId := ctx.Locals("userId").(string)
	motif := entity.Motif{
		Name:   req.Name,
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
