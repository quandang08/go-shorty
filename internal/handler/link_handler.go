package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quandang08/go-shorty/internal/model"
	"github.com/quandang08/go-shorty/internal/service"
)

// LinkHandler handles HTTP requests related to URL shortening.
type LinkHandler struct {
	Service service.LinkService
}

// NewLinkHandler creates a new LinkHandler instance.
func NewLinkHandler(svc service.LinkService) *LinkHandler {
	return &LinkHandler{Service: svc}
}

// CreateLink handles POST /api/v1/links.
// It validates the request, calls the service to create a short link,
// and returns the result or an error response.
func (h *LinkHandler) CreateLink(c *gin.Context) {
	var req model.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	response, err := h.Service.CreateShortLink(req.OriginalURL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The provided URL is invalid or malformed."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred."})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Redirect handles GET /:short_code.
// It resolves the short code, increments click count,
// and issues a 302 redirect to the original URL.
func (h *LinkHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("short_code")
	if shortCode == "" {
		c.Status(http.StatusNotFound)
		return
	}

	originalURL, err := h.Service.GetOriginalURL(shortCode)
	if err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short link not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

// GetLinkInfo handles GET /api/v1/links/:id.
// It retrieves metadata for a given short code without incrementing clicks.
func (h *LinkHandler) GetLinkInfo(c *gin.Context) {
	shortCode := c.Param("id")

	link, err := h.Service.GetLinkDetails(shortCode)
	if err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link details not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error retrieving link info"})
		return
	}
	c.JSON(http.StatusOK, link)
}

// ListLinks handles GET /api/v1/links.
// It returns a list of all short links as response DTOs.
func (h *LinkHandler) ListLinks(c *gin.Context) {
	links, err := h.Service.ListAllLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error retrieving link list"})
		return
	}
	c.JSON(http.StatusOK, links)
}
