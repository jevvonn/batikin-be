package bootstrap

import (
	"batikin-be/config"
	"batikin-be/internal/infra/postgresql"
	"fmt"

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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Your API is running")
	})

	return app.Listen(fmt.Sprintf(":%s", conf.AppPort))
}
