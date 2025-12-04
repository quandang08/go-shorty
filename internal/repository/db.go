package repository

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/quandang08/go-shorty/config"
	"github.com/quandang08/go-shorty/internal/model"
)

// InitDB khởi tạo kết nối GORM với PostgreSQL.
func InitDB(cfg *config.Config) *gorm.DB {
	// Mở kết nối
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Thiết lập Connection Pool (Cải thiện hiệu suất)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting DB instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// AutoMigrate: Tự động tạo bảng 'links' từ struct Link
	err = db.AutoMigrate(&model.Link{})
	if err != nil {
		log.Fatalf("Error running auto migration: %v", err)
	}

	log.Println("Database connection established and migration completed successfully.")
	return db
}
