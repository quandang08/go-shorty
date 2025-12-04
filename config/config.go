package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DSN         string
	ShortDomain string // Tên miền rút gọn (ví dụ: http://short.url)
}

// LoadConfig đọc các biến môi trường và trả về struct Config
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, loading config from environment variables.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		dbport := os.Getenv("DB_PORT")

		if host == "" {
			host = "localhost"
			user = "user"
			password = "password"
			dbname = "shorty_db"
			dbport = "5432"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
			host, user, password, dbname, dbport)
		log.Println("Using generated PostgreSQL DSN.")
	}

	// Domain rút gọn (quan trọng để trả về link đã rút gọn)
	shortDomain := os.Getenv("SHORT_DOMAIN")
	if shortDomain == "" {
		shortDomain = fmt.Sprintf("http://localhost:%s/", port)
	}

	return &Config{
		ServerPort:  port,
		DSN:         dsn,
		ShortDomain: shortDomain,
	}
}
