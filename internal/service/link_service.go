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
	// CreateShortLink generates a short URL for the given original URL.
	CreateShortLink(originalURL string) (*model.LinkResponse, error)

	// GetOriginalURL returns the original URL mapped to a short code
	// and increments its click count.
	GetOriginalURL(shortCode string) (string, error)

	// GetLinkDetails returns metadata for a short link without
	// incrementing the click count.
	GetLinkDetails(shortCode string) (*model.LinkResponse, error)

	// ListAllLinks returns all short links as response DTOs.
	ListAllLinks() ([]model.LinkResponse, error)
}

type linkServiceImpl struct {
	Repo repository.LinkRepository
	Cfg  *config.Config
}

// NewLinkService creates a new LinkService instance.
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

// CreateShortLink validates the URL, inserts a record to get the ID,
// generates a Base62 short code, updates the database, and returns metadata.
func (s *linkServiceImpl) CreateShortLink(originalURL string) (*model.LinkResponse, error) {
	if originalURL == "" {
		return nil, ErrInvalidURL
	}
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return nil, ErrInvalidURL
	}

	link := &model.Link{OriginalURL: originalURL}
	if err := s.Repo.Create(link); err != nil {
		return nil, ErrServiceUnavailable
	}

	link.ShortCode = util.EncodeToBase62(link.ID)
	if err := s.Repo.UpdateShortCode(link); err != nil {
		return nil, ErrServiceUnavailable
	}

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

// GetOriginalURL returns the original URL for a short code and
// increments the click counter.
func (s *linkServiceImpl) GetOriginalURL(shortCode string) (string, error) {
	if shortCode == "" {
		return "", ErrLinkNotFound
	}

	link, err := s.Repo.FindByShortCode(shortCode)
	if err != nil {
		return "", ErrServiceUnavailable
	}
	if link == nil {
		return "", ErrLinkNotFound
	}

	if err := s.Repo.IncrementClicks(shortCode); err != nil {
		return "", ErrServiceUnavailable
	}

	return link.OriginalURL, nil
}

// GetLinkDetails returns metadata for a short link without
// incrementing the click count.
func (s *linkServiceImpl) GetLinkDetails(shortCode string) (*model.LinkResponse, error) {
	link, err := s.Repo.FindByShortCode(shortCode)
	if err != nil {
		return nil, ErrServiceUnavailable
	}
	if link == nil {
		return nil, ErrLinkNotFound
	}

	return &model.LinkResponse{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		ClicksCount: link.ClicksCount,
		CreatedAt:   link.CreatedAt,
		ShortURL:    s.Cfg.ShortDomain + link.ShortCode,
	}, nil
}

// ListAllLinks retrieves all short links and maps them to response DTOs.
func (s *linkServiceImpl) ListAllLinks() ([]model.LinkResponse, error) {
	links, err := s.Repo.FindAll()
	if err != nil {
		return nil, ErrServiceUnavailable
	}

	responses := make([]model.LinkResponse, 0, len(links))
	shortDomain := s.Cfg.ShortDomain
	for _, link := range links {
		responses = append(responses, model.LinkResponse{
			ShortCode:   link.ShortCode,
			OriginalURL: link.OriginalURL,
			ClicksCount: link.ClicksCount,
			CreatedAt:   link.CreatedAt,
			ShortURL:    shortDomain + link.ShortCode,
		})
	}
	return responses, nil
}
