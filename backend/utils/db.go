package utils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	if err = createDefaultAdmin(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	adminTable := `
	CREATE TABLE IF NOT EXISTS admin (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	categoryTable := `
	CREATE TABLE IF NOT EXISTS category (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		sort INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	resourceTable := `
	CREATE TABLE IF NOT EXISTS resource (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER DEFAULT 0,
		title TEXT NOT NULL,
		cover TEXT DEFAULT '',
		description TEXT DEFAULT '',
		link TEXT NOT NULL,
		status INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(adminTable); err != nil {
		return err
	}

	if _, err := DB.Exec(categoryTable); err != nil {
		return err
	}

	if _, err := DB.Exec(resourceTable); err != nil {
		return err
	}

	// Add category_id column if it doesn't exist (for existing databases)
	DB.Exec("ALTER TABLE resource ADD COLUMN category_id INTEGER DEFAULT 0")

	// Create default categories if none exist
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM category").Scan(&count)
	if count == 0 {
		defaultCategories := []struct {
			name        string
			description string
			sort        int
		}{
			{"目标检测数据集", "目标检测相关数据集", 1},
			{"分类数据集", "图像分类相关数据集", 2},
			{"其他资源", "其他类型的资源", 99},
		}
		for _, c := range defaultCategories {
			DB.Exec("INSERT INTO category (name, description, sort) VALUES (?, ?, ?)",
				c.name, c.description, c.sort)
		}
		log.Println("Default categories created")
	}

	return nil
}

func createDefaultAdmin() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM admin").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("jiangcc8484"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = DB.Exec(
			"INSERT INTO admin (username, password) VALUES (?, ?)",
			"jeremyj",
			string(hashedPassword),
		)
		if err != nil {
			return err
		}
		log.Println("Default admin account created: jeremyj / jiangcc8484")
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
