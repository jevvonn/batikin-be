package bootstrap

import (
	"batikin-be/config"
	"batikin-be/internal/infra/midtrans"
	openaisdk "batikin-be/internal/infra/openai-sdk"
	"batikin-be/internal/infra/postgresql"
	"batikin-be/internal/infra/validator"
	"fmt"

	authHandler "batikin-be/internal/app/auth/interface/rest"
	motifHandler "batikin-be/internal/app/motif/interface/rest"
	orderHandler "batikin-be/internal/app/order/interface/rest"
	productHandler "batikin-be/internal/app/product/interface/rest"
	transactionHandler "batikin-be/internal/app/transaction/interface/rest"

	authUsecase "batikin-be/internal/app/auth/usecase"
	motifUsecase "batikin-be/internal/app/motif/usecase"
	orderUsecase "batikin-be/internal/app/order/usecase"
	productUsecase "batikin-be/internal/app/product/usecase"
	transactionUsecase "batikin-be/internal/app/transaction/usecase"

	motifRepository "batikin-be/internal/app/motif/repository"
	orderRepository "batikin-be/internal/app/order/repository"
	productRepository "batikin-be/internal/app/product/repository"
	transactionRepository "batikin-be/internal/app/transaction/repository"
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
	midtransClient := midtrans.NewMidtransSnapClient()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Your API is running")
	})

	apiRouter := app.Group("/api")

	userR := userRepository.NewUserPostgreSQL(db)
	motifR := motifRepository.NewMotifPostgreSQL(db)
	productR := productRepository.NewProductPostgreSQL(db)
	orderR := orderRepository.NewOrderPostgreSQL(db)
	transactionR := transactionRepository.NewTransactionPostgreSQL(db)

	authU := authUsecase.NewAuthUsecase(userR)
	motifU := motifUsecase.NewMotifUsecase(motifR, openAIClient)
	productU := productUsecase.NewProductUsecase(productR, motifR)
	orderU := orderUsecase.NewOrderUsecase(orderR, productR, midtransClient, transactionR)
	transactionU := transactionUsecase.NewTransactionUsecase(transactionR, orderR)

	authHandler.NewAuthHandler(apiRouter, authU, validator)
	motifHandler.NewMotifHandler(apiRouter, motifU, validator)
	productHandler.NewProductHandler(apiRouter, productU, validator)
	orderHandler.NewOrderHandler(apiRouter, orderU, validator)
	transactionHandler.NewTransactionHandler(apiRouter, transactionU)

	addr := fmt.Sprintf("localhost:%s", conf.AppPort)
	if conf.AppEnv == "production" {
		addr = fmt.Sprintf(":%s", conf.AppPort)
	}

	return app.Listen(addr)
}
