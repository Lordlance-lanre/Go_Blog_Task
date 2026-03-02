package main

import (
	database "Go_Blog_Task/Database"
	"os"
	"log"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"Go_Blog_Task/routes"
)

func main() {
	database.ConnectDB()
	err:=godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	port:=os.Getenv("PORT")
	app:=fiber.New()
	routes.Setup(app)
	app.Listen(":"+port)
}