package controller

import (
	"math/rand"

	"github.com/gofiber/fiber/v3"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(n int) string {
	// rand.seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// UploadImages godoc
// @Summary Upload an image
// @Description Uploads an image file and returns its URL.
// @Tags images
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/image/upload [post]
func UploadImages(c fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse body",
			"error":   err.Error(),
		})
	}
	files := form.File["image"]
	fileName := ""
	for _, file := range files {
		fileName = randomString(32) + "_" + file.Filename
		if err := c.SaveFile(file, "./uploads/"+fileName); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to save file",
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"url":     "http://localhost:8080/api/uploads/" + fileName,
		"file":    fileName,
	})
}
