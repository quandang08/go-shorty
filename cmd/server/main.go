package main

import (
	"log"

	"github.com/quandang08/go-shorty/config"
	"github.com/quandang08/go-shorty/internal/repository"
)

func main() {
	// Load cấu hình từ .env
	cfg := config.LoadConfig()

	// Khởi tạo Database và chạy Migration
	db := repository.InitDB(cfg)

	// Kiểm tra kết nối cuối cùng
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Server is ready to start...")
}
