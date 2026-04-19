package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"jz_web/models"
	"jz_web/utils"
)

// Public: Get all published resources
func GetResources(c *gin.Context) {
	rows, err := utils.DB.Query(
		"SELECT id, category_id, title, cover, description, link, status, created_at, updated_at FROM resource WHERE status = 1 ORDER BY category_id ASC, created_at DESC",
	)
	if err != nil {
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

// Public: Get single resource by ID
func GetResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	var r models.Resource
	err = utils.DB.QueryRow(
		"SELECT id, category_id, title, cover, description, link, status, created_at, updated_at FROM resource WHERE id = ? AND status = 1",
		id,
	).Scan(&r.ID, &r.CategoryID, &r.Title, &r.Cover, &r.Description, &r.Link, &r.Status, &r.CreatedAt, &r.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": r})
}

// Admin: Get all resources including unpublished
func GetAllResources(c *gin.Context) {
	categoryID := c.Query("category_id")

	var rows *sql.Rows
	var err error

	if categoryID != "" {
		rows, err = utils.DB.Query(
			"SELECT id, category_id, title, cover, description, link, status, created_at, updated_at FROM resource WHERE category_id = ? ORDER BY id ASC",
			categoryID,
		)
	} else {
		rows, err = utils.DB.Query(
			"SELECT id, category_id, title, cover, description, link, status, created_at, updated_at FROM resource ORDER BY category_id ASC, created_at DESC",
		)
	}

	if err != nil {
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

// Admin: Create new resource
func CreateResource(c *gin.Context) {
	var req models.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Auto-generate next sequential ID
	var maxID int64
	utils.DB.QueryRow("SELECT COALESCE(MAX(id), 0) FROM resource").Scan(&maxID)
	newID := maxID + 1

	categoryID := req.CategoryID
	if categoryID == 0 {
		categoryID = 1 // Default to first category
	}

	result, err := utils.DB.Exec(
		"INSERT INTO resource (id, category_id, title, cover, description, link, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		newID, categoryID, req.Title, req.Cover, req.Description, req.Link, req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resource"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"message": "Resource created",
		"data": gin.H{"id": id},
	})
}

// Admin: Update resource
func UpdateResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	var req models.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	categoryID := req.CategoryID
	if categoryID == 0 {
		categoryID = 1
	}

	_, err = utils.DB.Exec(
		"UPDATE resource SET category_id = ?, title = ?, cover = ?, description = ?, link = ?, status = ?, updated_at = ? WHERE id = ?",
		categoryID, req.Title, req.Cover, req.Description, req.Link, req.Status, time.Now(), id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource updated"})
}

// Admin: Delete resource
func DeleteResource(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	// Start transaction
	tx, err := utils.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource"})
		return
	}

	// Delete the resource
	_, err = tx.Exec("DELETE FROM resource WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource"})
		return
	}

	// Renumber resources with ID > deletedID (shift down by 1)
	_, err = tx.Exec("UPDATE resource SET id = id - 1 WHERE id > ?", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to renumber resources"})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted"})
}

// Admin: Toggle resource status
func ToggleResourceStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	var currentStatus int
	err = utils.DB.QueryRow("SELECT status FROM resource WHERE id = ?", id).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	newStatus := 0
	if currentStatus == 0 {
		newStatus = 1
	}

	_, err = utils.DB.Exec("UPDATE resource SET status = ?, updated_at = ? WHERE id = ?", newStatus, time.Now(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status toggled", "status": newStatus})
}
