package bootstrap

import (
	"batikin-be/config"
	openaisdk "batikin-be/internal/infra/openai-sdk"
	"batikin-be/internal/infra/postgresql"
	"batikin-be/internal/infra/validator"
	"fmt"

	authHandler "batikin-be/internal/app/auth/interface/rest"
	motifHandler "batikin-be/internal/app/motif/interface/rest"

	authUsecase "batikin-be/internal/app/auth/usecase"
	motifUsecase "batikin-be/internal/app/motif/usecase"

	motifRepository "batikin-be/internal/app/motif/repository"
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
	openAIClient := openaisdk.NewOpenAIClient()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Your API is running")
	})

	apiRouter := app.Group("/api")

	userR := userRepository.NewUserPostgreSQL(db)
	motifR := motifRepository.NewMotifPostgreSQL(db)

	authU := authUsecase.NewAuthUsecase(userR)
	motifU := motifUsecase.NewMotifUsecase(motifR, openAIClient)

	authHandler.NewAuthHandler(apiRouter, authU, validator)
	motifHandler.NewMotifHandler(apiRouter, motifU, validator)

	addr := fmt.Sprintf("localhost:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		addr = fmt.Sprintf(":%s", conf.AppPort)
	}

	return app.Listen(addr)
}
