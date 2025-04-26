package middleware

import (
	"batikin-be/internal/infra/jwt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Authenticated(ctx *fiber.Ctx) error {
	headers := ctx.Get("Authorization")

	if headers == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})

	}

	tokenString := strings.Replace(headers, "Bearer ", "", 1)
	claims, err := jwt.ParseAuthToken(tokenString)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	exp := claims["exp"].(float64)
	expiredDate := time.Unix(int64(exp), 0)

	if expiredDate.Before(time.Now()) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token expired",
		})
	}

	ctx.Locals("userId", claims["sub"])
	ctx.Locals("email", claims["email"])
	ctx.Locals("name", claims["name"])

	return ctx.Next()
}
