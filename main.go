package main

import (
	database "Go_Blog_Task/Database"
	"Go_Blog_Task/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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

	// Enable CORS middleware before routes setup
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	routes.Setup(app)
	app.Listen(":" + port)
}
