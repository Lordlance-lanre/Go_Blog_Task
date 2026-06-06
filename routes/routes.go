package routes

import (
	"Go_Blog_Task/Database/controller"
	"Go_Blog_Task/middleware"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func Setup(app *fiber.App) {
	auth := app.Group("/api/auth")
	auth.Post("/register", controller.Register)

	// Rate limiter for login: max 5 requests per minute
	loginLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many login attempts. Please try again after a minute.",
			})
		},
	})
	auth.Post("/login", loginLimiter, controller.Login)

	app.Use(middleware.AuthGuard) //anywhere this middleware is placed signifies, the routes that should be authenticated.

	blog := app.Group("/api/blog")
	blog.Post("/create", controller.CreatePost)
	blog.Get("all-posts", controller.GetAllPost)
	blog.Get("/user-posts", controller.GetUniquePosts)
	blog.Get("/:id", controller.GetSinglePost)
	blog.Put("/update-post/:id", controller.UpdateBlog)
	blog.Delete("/delete-post/:id", controller.DeleteBlogPost)

	image := app.Group("/api/image")
	image.Post("/upload", controller.UploadImages)
	image.Get("/uploads/*", static.New("./uploads"))

	// auth.Get("/user", controller.User)
}
