package service

import (
	"errors"
	"net/url"

	"github.com/quandang08/go-shorty/config"
	"github.com/quandang08/go-shorty/internal/model"
	"github.com/quandang08/go-shorty/internal/repository"
	"github.com/quandang08/go-shorty/internal/util"
)

// LinkService defines the business logic for creating and resolving short URLs.
type LinkService interface {
	// CreateShortLink generates a short code for the given URL and returns link metadata.
	CreateShortLink(originalURL string) (*model.LinkResponse, error)

	// GetOriginalURL returns the long URL mapped to the given short code.
	GetOriginalURL(shortCode string) (string, error)
}

// linkServiceImpl implements LinkService.
type linkServiceImpl struct {
	Repo repository.LinkRepository
	Cfg  *config.Config
}

// NewLinkService initializes and returns a LinkService instance.
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

// CreateShortLink validates the URL, generates a Base62 short code from the DB ID,
// persists it, and returns the full response payload.
func (s *linkServiceImpl) CreateShortLink(originalURL string) (*model.LinkResponse, error) {
	// Validate input
	if originalURL == "" {
		return nil, ErrInvalidURL
	}
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return nil, ErrInvalidURL
	}

	// Insert initial record to obtain auto-increment ID
	link := &model.Link{OriginalURL: originalURL}
	if err := s.Repo.Save(link); err != nil {
		return nil, ErrServiceUnavailable
	}

	// Generate short code from ID
	shortCode := util.EncodeToBase62(link.ID)
	link.ShortCode = shortCode

	// Save short code back into database
	if err := s.Repo.Save(link); err != nil {
		return nil, ErrServiceUnavailable
	}

	// Build short URL
	shortURL := s.Cfg.ShortDomain
	if len(shortURL) > 0 && shortURL[len(shortURL)-1] != '/' {
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

// GetOriginalURL retrieves the original URL mapped to a short code
// and increments the click counter.
func (s *linkServiceImpl) GetOriginalURL(shortCode string) (string, error) {
	if shortCode == "" {
		return "", ErrLinkNotFound
	}

	// Lookup link
	link, err := s.Repo.FindByShortCode(shortCode)
	if err != nil {
		return "", ErrServiceUnavailable
	}
	if link == nil {
		return "", ErrLinkNotFound
	}

	// Increment click count (repository should ensure atomic update)
	if err := s.Repo.IncrementClicks(shortCode); err != nil {
		return "", ErrServiceUnavailable
	}

	return link.OriginalURL, nil
}
