package controller

import (
	database "Go_Blog_Task/Database"
	"Go_Blog_Task/models"
	"Go_Blog_Task/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// CreatePost godoc
// @Summary Create a blog post
// @Description Creates a new blog post.
// @Tags blogs
// @Accept json
// @Produce json
// @Param blog body models.Blog true "Blog post data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/blog/create [post]
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
// GetAllPost godoc
// @Summary Get all blog posts
// @Description Returns a paginated list of blog posts.
// @Tags blogs
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/blog/all-posts [get]
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
// GetSinglePost godoc
// @Summary Get a single blog post
// @Description Returns details of a specific blog post.
// @Tags blogs
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/blog/:id [get]

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

// UpdateBlog godoc
// @Summary Update a blog post
// @Description Updates an existing blog post by its ID.
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path int true "Blog post ID"
// @Param blog body models.Blog true "Updated blog post data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/blog/update-post/{id} [put]
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
// GetUniquePosts godoc
// @Summary Get the current user's posts
// @Description Returns all blog posts belonging to the authenticated user.
// @Tags blogs
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/blog/user-posts [get]
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

// DeleteBlogPost godoc
// @Summary Delete a blog post
// @Description Deletes a blog post by its ID.
// @Tags blogs
// @Produce json
// @Param id path int true "Blog post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/blog/delete-post/{id} [delete]
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
