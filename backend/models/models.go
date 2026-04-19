package models

import "time"

type Admin struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID          int64     `json:"id"`
	Name       string    `json:"name"`
	Description string    `json:"description"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
}

type Resource struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Title       string    `json:"title"`
	Cover       string    `json:"cover"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Admin Admin  `json:"admin"`
}

type CreateCategoryRequest struct {
	Name       string `json:"name" binding:"required"`
	Description string `json:"description"`
	Sort       int    `json:"sort"`
}

type UpdateCategoryRequest struct {
	Name       string `json:"name"`
	Description string `json:"description"`
	Sort       int    `json:"sort"`
}

type CreateResourceRequest struct {
	CategoryID  int64  `json:"category_id"`
	Title       string `json:"title" binding:"required"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Link        string `json:"link" binding:"required"`
	Status      int    `json:"status"`
}

type UpdateResourceRequest struct {
	CategoryID  int64  `json:"category_id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Status      int    `json:"status"`
}