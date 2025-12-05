package service

import (
	"errors"

	"github.com/quandang08/go-shorty/config"
	"github.com/quandang08/go-shorty/internal/model"
	"github.com/quandang08/go-shorty/internal/repository"
)

// LinkService defines the business logic for creating and resolving short URLs.
type LinkService interface {
	CreateShortLink(originalURL string) (*model.LinkResponse, error)
	GetOriginalURL(shortCode string) (string, error)
}

type linkServiceImpl struct {
	Repo repository.LinkRepository
	Cfg  *config.Config
}

func NewLinkService(repo repository.LinkRepository, cfg *config.Config) LinkService {
	return &linkServiceImpl{
		Repo: repo,
		Cfg:  cfg,
	}
}

// Business errors returned by the service layer.
var (
	ErrInvalidURL         = errors.New("URL is invalid or empty")
	ErrLinkNotFound       = errors.New("short link not found")
	ErrConflict           = errors.New("short code conflict occurred, please retry")
	ErrServiceUnavailable = errors.New("internal service error")
)
