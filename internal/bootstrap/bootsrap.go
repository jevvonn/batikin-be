package bootstrap

import (
	"batikin-be/config"
	"batikin-be/internal/infra/postgresql"
	"batikin-be/internal/infra/validator"
	"fmt"

	authHandler "batikin-be/internal/app/auth/interface/rest"
	authUsecase "batikin-be/internal/app/auth/usecase"
	userRepository "batikin-be/internal/app/user/repository"

	"github.com/gofiber/fiber/v2"
)

func Start() error {
	app := fiber.New()
	conf := config.New()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		conf.DbHost,
		conf.DbPort,
		conf.DbUser,
		conf.DbPassword,
		conf.DbName,
	)
	db, err := postgresql.New(dsn)
	if err != nil {
		panic(err)
	}

	CommandHandler(db)

	validator := validator.NewValidator()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Your API is running")
	})

	apiRouter := app.Group("/api")

	userR := userRepository.NewUserPostgreSQL(db)

	authU := authUsecase.NewAuthUsecase(userR)

	authHandler.NewAuthHandler(apiRouter, authU, validator)

	return app.Listen(fmt.Sprintf("localhost:%s", conf.AppPort))
}
