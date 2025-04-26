package usecase

import (
	motifRepo "batikin-be/internal/app/motif/repository"
	"batikin-be/internal/app/product/repository"
	"batikin-be/internal/constant"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/domain/entity"
	"batikin-be/internal/helper"
	"batikin-be/internal/infra/generation"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"gorm.io/gorm"
)

type ProductUsecaseItf interface {
	GetAll() ([]entity.Product, error)
	GetSpecific(ctx *fiber.Ctx) (entity.Product, error)
	CreateFromMotif(ctx *fiber.Ctx, req dto.CreateFromMotifProductRequest) (entity.Product, error)
}

type ProductUsecase struct {
	productRepo  repository.ProductPostgreSQLItf
	motifRepo    motifRepo.MotifPostgreSQLItf
	openAIClient openai.Client
}

func NewProductUsecase(
	productRepo repository.ProductPostgreSQLItf,
	motifRepo motifRepo.MotifPostgreSQLItf,
	openAIClient openai.Client,
) ProductUsecaseItf {
	return &ProductUsecase{productRepo, motifRepo, openAIClient}
}

func (u *ProductUsecase) GetAll() ([]entity.Product, error) {
	return u.productRepo.GetAll()
}

func (u *ProductUsecase) GetSpecific(ctx *fiber.Ctx) (entity.Product, error) {
	param := ctx.Params("id")
	productId, err := uuid.Parse(param)
	if err != nil {
		return entity.Product{}, err
	}

	product := entity.Product{ID: productId}

	return u.productRepo.GetSpecific(product)
}

func (u *ProductUsecase) CreateFromMotif(ctx *fiber.Ctx, req dto.CreateFromMotifProductRequest) (entity.Product, error) {
	param := ctx.Params("motifId")
	motifId, err := uuid.Parse(param)

	if err != nil {
		return entity.Product{}, err
	}

	motif, err := u.motifRepo.GetSpecific(entity.Motif{ID: motifId})
	if err != nil {
		return entity.Product{}, err
	}

	typeCloth := constant.CLOTH_TYPE[req.ClothType]

	checkProduct, err := u.productRepo.GetSpecific(entity.Product{
		MotifID:   motif.ID,
		ClothType: typeCloth,
	})

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Product{}, err
		}
	}

	if checkProduct.ID != uuid.Nil {
		return checkProduct, nil
	}

	_, url, err := helper.FetchAndSaveImage(motif.ImageURL, "tmp")
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return entity.Product{}, err
	}

	prompt := fmt.Sprintf("buatkan menjadi design mockup %s tanpa orang untuk keperluan AR test-on filter", typeCloth)
	editedURL, err := generation.GenerateImageEdit(url, prompt)
	if err != nil {
		return entity.Product{}, err
	}

	fmt.Println("Image URL:", editedURL)
	productId := uuid.New()
	sizes := []entity.ProductSizeVariant{}
	for _, size := range constant.CLOTH_SIZE {
		variantId := uuid.New()
		sizes = append(sizes, entity.ProductSizeVariant{
			ID:        variantId,
			Size:      size.Size,
			Price:     size.Price,
			ProductID: productId,
		})
	}

	product := &entity.Product{
		ID:        productId,
		Name:      typeCloth + " " + motif.Prompt,
		ImageURL:  editedURL,
		Sizes:     sizes,
		ClothType: typeCloth,
		MotifID:   motif.ID,
	}

	if err := u.productRepo.Create(product); err != nil {
		return entity.Product{}, err
	}

	response, err := u.productRepo.GetSpecific(entity.Product{ID: productId})
	if err != nil {
		return entity.Product{}, err
	}

	return response, nil
}
