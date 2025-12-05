package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quandang08/go-shorty/internal/model"
	"github.com/quandang08/go-shorty/internal/service"
)

// LinkHandler handles HTTP requests related to URL shortening.
type LinkHandler struct {
	Service service.LinkService
}

// NewLinkHandler returns a new instance of LinkHandler.
func NewLinkHandler(svc service.LinkService) *LinkHandler {
	return &LinkHandler{Service: svc}
}

// CreateLink handles POST /api/v1/links
// It validates the request body, delegates business logic to the service layer,
// and returns the shortened URL metadata.
func (h *LinkHandler) CreateLink(c *gin.Context) {
	var req model.CreateLinkRequest

	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input format or missing fields",
		})
		return
	}

	// Create short link using service layer
	response, err := h.Service.CreateShortLink(req.OriginalURL)
	if err != nil {
		if err == service.ErrInvalidURL {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "The provided URL is invalid or malformed.",
			})
			return
		}

		// Unexpected internal error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error occurred.",
		})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, response)
}

// Redirect handles GET /:short_code
// It looks up the original URL and issues an HTTP redirect.
func (h *LinkHandler) Redirect(c *gin.Context) {
	// Logic sẽ được bổ sung
}

// GetLinkInfo handles GET /api/v1/links/:id
func (h *LinkHandler) GetLinkInfo(c *gin.Context) {
	// Logic sẽ được bổ sung
}

// ListLinks handles GET /api/v1/links
func (h *LinkHandler) ListLinks(c *gin.Context) {
	// Logic sẽ được bổ sung
}
