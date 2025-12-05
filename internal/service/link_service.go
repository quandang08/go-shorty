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
	GetLinkDetails(shortCode string) (*model.LinkResponse, error) // XEM CHI TIẾT
	ListAllLinks() ([]model.LinkResponse, error)                  // LIỆT KÊ
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

// GetLinkDetails Implement GetLinkDetails (chỉ gọi FindByShortCode, không tăng click)
func (s *linkServiceImpl) GetLinkDetails(shortCode string) (*model.LinkResponse, error) {
	link, err := s.Repo.FindByShortCode(shortCode)
	if err != nil {
		return nil, ErrServiceUnavailable
	}
	if link == nil {
		return nil, ErrLinkNotFound
	}
	// Map Entity sang Response DTO (Sử dụng code này để tránh trùng lặp)
	return &model.LinkResponse{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		ClicksCount: link.ClicksCount,
		CreatedAt:   link.CreatedAt,
		ShortURL:    s.Cfg.ShortDomain + link.ShortCode,
	}, nil
}

func (s *linkServiceImpl) ListAllLinks() ([]model.LinkResponse, error) {
	links, err := s.Repo.FindAll()
	if err != nil {
		return nil, ErrServiceUnavailable
	}
	// ... (logic map links sang responses DTO)
	var responses []model.LinkResponse
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
