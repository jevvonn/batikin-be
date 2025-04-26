package usecase

import (
	userRepo "batikin-be/internal/app/user/repository"
	"batikin-be/internal/domain/dto"
	"batikin-be/internal/domain/entity"
	"batikin-be/internal/helper"
	"batikin-be/internal/infra/jwt"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	Register(ctx *fiber.Ctx, req dto.RegisterRequest) error
	Login(ctx *fiber.Ctx, req dto.LoginRequest) (dto.LoginResponse, error)
	Session(ctx *fiber.Ctx) (dto.SessionResponse, error)
}

type AuthUsecase struct {
	userRepo userRepo.UserPostgreSQLItf
}

func NewAuthUsecase(
	userRepo userRepo.UserPostgreSQLItf,
) AuthUsecaseItf {
	return &AuthUsecase{userRepo}
}

func (u *AuthUsecase) Register(ctx *fiber.Ctx, req dto.RegisterRequest) error {
	user, err := u.userRepo.GetSpecificUser(entity.User{
		Email: req.Email,
	})

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if user.ID != uuid.Nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user = entity.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}
	err = u.userRepo.CreateUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) Login(ctx *fiber.Ctx, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := u.userRepo.GetSpecificUser(entity.User{
		Email: req.Email,
	})

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.LoginResponse{}, err
	}

	if user.ID == uuid.Nil {
		return dto.LoginResponse{}, errors.New("email or password is incorrect")
	}

	if !helper.VerifyPassword(req.Password, user.Password) {
		return dto.LoginResponse{}, errors.New("email or password is incorrect")
	}

	token, err := jwt.CreateAuthToken(user.ID.String(), user.Email, user.Name)

	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		UserId: user.ID.String(),
		Token:  token,
	}, nil
}

func (u *AuthUsecase) Session(ctx *fiber.Ctx) (dto.SessionResponse, error) {
	userId := ctx.Locals("userId").(string)

	uuidUser, err := uuid.Parse(userId)
	if err != nil {
		return dto.SessionResponse{}, err
	}

	user, err := u.userRepo.GetSpecificUser(entity.User{
		ID: uuidUser,
	})
	if err != nil {
		return dto.SessionResponse{}, err
	}

	return dto.SessionResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
