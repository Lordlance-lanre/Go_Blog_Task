package middleware

import (
	"Go_Blog_Task/utils"
	"github.com/gofiber/fiber/v3"
)

func AuthGuard(c fiber.Ctx) error {
	cookie:= c.Cookies("jwt")

	if _, err := utils.ParseJWT(cookie); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized Access",
		})
	}
	return c.Next()
	
}