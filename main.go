package main

import (
	database "Go_Blog_Task/Database"
	"Go_Blog_Task/routes"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	database.ConnectDB()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	port := os.Getenv("PORT")
	app := fiber.New()

	// Dynamic CORS middleware: allow origins from env or Vercel deployments
	app.Use(func(c fiber.Ctx) error {
		origin := c.Get("Origin")

		// Build allowed origins list from environment variable `ALLOW_ORIGINS`
		// (comma-separated) and include the current Vercel URL when available.
		allowed := []string{"http://localhost:3000", "http://localhost:5173"}
		if env := os.Getenv("ALLOW_ORIGINS"); env != "" {
			for _, o := range strings.Split(env, ",") {
				o = strings.TrimSpace(o)
				if o != "" {
					allowed = append(allowed, o)
				}
			}
		}
		if v := os.Getenv("VERCEL_URL"); v != "" {
			allowed = append(allowed, "https://"+v)
		}

		allow := false
		if origin != "" {
			for _, a := range allowed {
				if origin == a || origin == a+"/" {
					allow = true
					break
				}
			}
			// Accept any vercel.app subdomain (preview deployments)
			if !allow && strings.HasSuffix(origin, ".vercel.app") {
				allow = true
			}
		}

		if allow {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Expose-Headers", "Set-Cookie")
		}

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	})

	routes.Setup(app)
	app.Listen(":" + port)
}
