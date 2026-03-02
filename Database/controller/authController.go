package controller

import (
	"Go_Blog_Task/Database"
	"Go_Blog_Task/models"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func validateEmail(email string) bool {
	regisMail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regisMail.MatchString(email)
}

func Register(c fiber.Ctx) error {
	var data map[string]interface{}

	var userData models.User

	if err := c.Bind().Body(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	//check if password is less than 6 characters

	if len(data["password"].(string)) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 6 characters or more",
		})
	}
	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email format",
		})
	}

	//check for existing email

	database.DB.Where("email = ?", strings.TrimSpace(data["email"].(string))).First(&userData)

	if userData.Id != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}
	
	userData = models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		Phone:     data["phone"].(string),
		// Password: data["password"].(string),
	}
	
	userData.SetPassword(data["password"].(string))
	err := database.DB.Create(&userData).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register user",
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    userData,
	})
}
