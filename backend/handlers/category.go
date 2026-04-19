package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"jz_web/models"
	"jz_web/utils"
)

// Get all categories (public)
func GetCategories(c *gin.Context) {
	rows, err := utils.DB.Query("SELECT id, name, description, sort, created_at FROM category ORDER BY sort ASC, id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.Sort, &cat.CreatedAt); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// Create category (admin only)
func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := utils.DB.Exec(
		"INSERT INTO category (name, description, sort) VALUES (?, ?, ?)",
		req.Name, req.Description, req.Sort,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "Category created", "data": gin.H{"id": id}})
}

// Update category (admin only)
func UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err = utils.DB.Exec(
		"UPDATE category SET name = ?, description = ?, sort = ? WHERE id = ?",
		req.Name, req.Description, req.Sort, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated"})
}

// Delete category (admin only)
func DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Set resources in this category to category_id = 0
	utils.DB.Exec("UPDATE resource SET category_id = 0 WHERE category_id = ?", id)

	_, err = utils.DB.Exec("DELETE FROM category WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

// Get resources by category (public)
func GetResourcesByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("categoryId"), 10, 64)
	if err != nil {
		categoryID = 0
	}

	keyword := c.Query("keyword")

	var rows *sql.Rows
	var queryErr error

	if keyword != "" {
		if categoryID > 0 {
			rows, queryErr = utils.DB.Query(
				`SELECT id, category_id, title, cover, description, link, status, created_at, updated_at
				FROM resource WHERE status = 1 AND category_id = ? AND (title LIKE ? OR description LIKE ?)
				ORDER BY created_at DESC`,
				categoryID, "%"+keyword+"%", "%"+keyword+"%",
			)
		} else {
			rows, queryErr = utils.DB.Query(
				`SELECT id, category_id, title, cover, description, link, status, created_at, updated_at
				FROM resource WHERE status = 1 AND (title LIKE ? OR description LIKE ?)
				ORDER BY created_at DESC`,
				"%"+keyword+"%", "%"+keyword+"%",
			)
		}
	} else {
		// No keyword filter
		if categoryID > 0 {
			rows, queryErr = utils.DB.Query(
				`SELECT id, category_id, title, cover, description, link, status, created_at, updated_at
				FROM resource
				WHERE status = 1 AND category_id = ?
				ORDER BY created_at DESC`,
				categoryID,
			)
		} else {
			rows, queryErr = utils.DB.Query(
				`SELECT id, category_id, title, cover, description, link, status, created_at, updated_at
				FROM resource
				WHERE status = 1
				ORDER BY created_at DESC`,
			)
		}
	}

	if queryErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}
	defer rows.Close()

	resources := []models.Resource{}
	for rows.Next() {
		var r models.Resource
		if err := rows.Scan(&r.ID, &r.CategoryID, &r.Title, &r.Cover, &r.Description, &r.Link, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			continue
		}
		resources = append(resources, r)
	}

	c.JSON(http.StatusOK, gin.H{"data": resources})
}

// Search resources (public)
func SearchResources(c *gin.Context) {
	keyword := c.Query("keyword")
	categoryID, _ := strconv.ParseInt(c.Query("category_id"), 10, 64)

	var rows *sql.Rows
	var err error

	if keyword != "" && categoryID > 0 {
		rows, err = utils.DB.Query(
			`SELECT id, category_id, title, cover, description, link, status, created_at, updated_at
			FROM resource
			WHERE status = 1 AND category_id = ? AND (title LIKE ? OR description LIKE ?)
			ORDER BY created_at DESC`,
			categoryID, "%"+keyword+"%", "%"+keyword+"%",
		)
	} else if keyword != "" {
		rows, err = utils.DB.Query(
			`SELECT id, category_id, title, cover, description, link, status, created_at, updated_at
			FROM resource
			WHERE status = 1 AND (title LIKE ? OR description LIKE ?)
			ORDER BY created_at DESC`,
			"%"+keyword+"%", "%"+keyword+"%",
		)
	} else if categoryID > 0 {
		rows, err = utils.DB.Query(
			`SELECT id, category_id, title, cover, description, link, status, created_at, updated_at
			FROM resource
			WHERE status = 1 AND category_id = ?
			ORDER BY created_at DESC`,
			categoryID,
		)
	} else {
		rows, err = utils.DB.Query(
			`SELECT id, category_id, title, cover, description, link, status, created_at, updated_at
			FROM resource
			WHERE status = 1
			ORDER BY created_at DESC`,
		)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search resources"})
		return
	}
	defer rows.Close()

	resources := []models.Resource{}
	for rows.Next() {
		var r models.Resource
		if err := rows.Scan(&r.ID, &r.CategoryID, &r.Title, &r.Cover, &r.Description, &r.Link, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			continue
		}
		resources = append(resources, r)
	}

	c.JSON(http.StatusOK, gin.H{"data": resources})
}