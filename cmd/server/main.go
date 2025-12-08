package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/quandang08/go-shorty/config"
	"github.com/quandang08/go-shorty/internal/handler"
	"github.com/quandang08/go-shorty/internal/repository"
	"github.com/quandang08/go-shorty/internal/service"
)

func main() {
	// Load configuration and initialize database.
	cfg := config.LoadConfig()
	db := repository.InitDB(cfg)

	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	// Setup dependencies.
	linkRepo := repository.NewLinkRepository(db)
	linkService := service.NewLinkService(linkRepo, cfg)
	linkHandler := handler.NewLinkHandler(linkService)

	// Initialize Gin router.
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `
		<h1>Welcome to GoShorty!</h1>
		<p>Enter a URL to shorten:</p>
		<form action="/api/v1/links" method="POST">
			<input type="url" name="original_url" placeholder="https://example.com" required style="width:300px">
			<button type="submit">Shorten</button>
		</form>
		<p>Or visit <a href="/api/v1/links">/api/v1/links</a> to see the list of short URLs (GET).</p>
	`)
	})

	// Public redirect route: GET /:short_code
	router.GET("/:short_code", linkHandler.Redirect)

	// API v1 routes.
	v1 := router.Group("/api/v1")
	{
		v1.POST("/links", linkHandler.CreateLink)
		v1.GET("/links/:id", linkHandler.GetLinkInfo)
		v1.GET("/links", linkHandler.ListLinks)
	}

	// Start server.
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
