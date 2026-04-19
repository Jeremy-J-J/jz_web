package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"jz_web/handlers"
	"jz_web/middleware"
	"jz_web/utils"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/workspace/pro/jz_web/data/db/app.db"
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal("Failed to create db directory:", err)
	}

	if err := utils.InitDB(dbPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer utils.CloseDB()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Static("/uploads", "/workspace/pro/jz_web/data/uploads")

	api := r.Group("/api")
	{
		api.POST("/admin/login", handlers.Login)

		api.GET("/categories", handlers.GetCategories)
		api.GET("/resources", handlers.GetResources)
		api.GET("/resources/:id", handlers.GetResource)
		api.GET("/search", handlers.SearchResources)
		api.GET("/stats", handlers.GetStats)
		api.POST("/visit", handlers.RecordVisit)

		admin := api.Group("/admin")
		admin.Use(middleware.AuthRequired())
		{
			admin.GET("/categories", handlers.GetCategories)
			admin.POST("/categories", handlers.CreateCategory)
			admin.PUT("/categories/:id", handlers.UpdateCategory)
			admin.DELETE("/categories/:id", handlers.DeleteCategory)

			admin.GET("/resources", handlers.GetAllResources)
			admin.POST("/resources", handlers.CreateResource)
			admin.PUT("/resources/:id", handlers.UpdateResource)
			admin.DELETE("/resources/:id", handlers.DeleteResource)
			admin.PUT("/resources/:id/toggle", handlers.ToggleResourceStatus)
			admin.POST("/upload", handlers.UploadCover)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
