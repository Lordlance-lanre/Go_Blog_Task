package routes

import(
	"Go_Blog_Task/Database/controller"
	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App) {
	auth := app.Group("/api/auth")
	auth.Post("/register", controller.Register)
}