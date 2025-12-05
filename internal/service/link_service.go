package service

import (
	"errors"
	"net/url"

	"github.com/quandang08/go-shorty/config"
	"github.com/quandang08/go-shorty/internal/model"
	"github.com/quandang08/go-shorty/internal/repository"
	"github.com/quandang08/go-shorty/internal/util"
)

// LinkService defines the URL-shortening business logic.
type LinkService interface {
	CreateShortLink(originalURL string) (*model.LinkResponse, error)
	GetOriginalURL(shortCode string) (string, error)
}

type linkServiceImpl struct {
	Repo repository.LinkRepository
	Cfg  *config.Config
}

func NewLinkService(repo repository.LinkRepository, cfg *config.Config) LinkService {
	return &linkServiceImpl{Repo: repo, Cfg: cfg}
}

// Business-level errors.
var (
	ErrInvalidURL         = errors.New("invalid URL")
	ErrLinkNotFound       = errors.New("short link not found")
	ErrConflict           = errors.New("short code conflict")
	ErrServiceUnavailable = errors.New("service unavailable")
)

// CreateShortLink validates the URL, stores the record to obtain ID,
// generates a Base62 short code, updates DB, and returns metadata.
func (s *linkServiceImpl) CreateShortLink(originalURL string) (*model.LinkResponse, error) {
	// Validate URL
	if originalURL == "" {
		return nil, ErrInvalidURL
	}
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return nil, ErrInvalidURL
	}

	// Insert to obtain auto-increment ID
	link := &model.Link{OriginalURL: originalURL}
	if err := s.Repo.Create(link); err != nil {
		return nil, ErrServiceUnavailable
	}

	// Generate short code from ID
	link.ShortCode = util.EncodeToBase62(link.ID)

	// Update DB with generated short code
	if err := s.Repo.UpdateShortCode(link); err != nil {
		return nil, ErrServiceUnavailable
	}

	// Build final short URL
	shortURL := s.Cfg.ShortDomain
	if shortURL != "" && shortURL[len(shortURL)-1] != '/' {
		shortURL += "/"
	}

	return &model.LinkResponse{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		ClicksCount: link.ClicksCount,
		CreatedAt:   link.CreatedAt,
		ShortURL:    shortURL + link.ShortCode,
	}, nil
}

// GetOriginalURL finds the original URL mapped to the short code
// and increments click count.
func (s *linkServiceImpl) GetOriginalURL(shortCode string) (string, error) {
	if shortCode == "" {
		return "", ErrLinkNotFound
	}

	// Lookup
	link, err := s.Repo.FindByShortCode(shortCode)
	if err != nil {
		return "", ErrServiceUnavailable
	}
	if link == nil {
		return "", ErrLinkNotFound
	}

	// Increment counter
	if err := s.Repo.IncrementClicks(shortCode); err != nil {
		return "", ErrServiceUnavailable
	}

	return link.OriginalURL, nil
}
