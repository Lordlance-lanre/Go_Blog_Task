package controller

import (
	database "Go_Blog_Task/Database"
	"Go_Blog_Task/models"
	"Go_Blog_Task/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func CreatePost(c fiber.Ctx) error {
	var blog models.Blog
	if err := c.Bind().Body(&blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse body",
			"error":   err.Error(),
		})
	}

	if err := database.DB.Create(&blog).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create post",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post created successfully",
		"blog":    blog,
	})

}

func GetAllPost(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 5
	offset := (page - 1) * limit
	var total int64
	var getBlog []models.Blog
	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getBlog)
	database.DB.Count(&total)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Posts retrieved successfully",
		"blog":    getBlog,
		"meta": fiber.Map{
			"page":      page,
			"limit":     limit,
			"offset":    offset,
			"total":     total,
			"last_page": int64(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

func GetSinglePost(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var getBlog models.Blog
	database.DB.Where("id=?", id).Preload("User").First(&getBlog)
	if getBlog.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post retrieved successfully",
		"blog":    getBlog,
	})
}

func UpdateBlog(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	blog := models.Blog{
		Id: uint(id),
	}
	if err := c.Bind().Body(&blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse body",
			"error":   err.Error(),
		})
	}
	database.DB.Model(&blog).Updates(&blog)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post updated successfully",
		"blog": fiber.Map{
			"id":          blog.Id,
			"title":       blog.Title,
			"description": blog.Description,
			"image_url":   blog.ImageURL,
			"user_id":     blog.UserID,
		},
	})

}

func GetUniquePosts(c fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, _ := utils.ParseJWT(cookie)
	var blogs []models.Blog
	database.DB.Where("user_id=?", id).Preload("User").Find(&blogs)
	if len(blogs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Posts retrieved successfully",
		"blog":    blogs,
	})

}

func DeleteBlogPost(c fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	blog := models.Blog{
		Id: uint(id),
	}
	if blog.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not found",
		})
	}
	result := database.DB.Delete(&blog)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete post",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post deleted successfully",
	})
}
