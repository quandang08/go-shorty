package service

import (
	"errors"
	"github.com/quandang08/go-shorty/config"
	"github.com/quandang08/go-shorty/internal/model"
	"github.com/quandang08/go-shorty/internal/repository"
)

// LinkService định nghĩa các phương thức nghiệp vụ
type LinkService interface {
	CreateShortLink(originalURL string) (*model.LinkResponse, error)
	GetOriginalURL(shortCode string) (string, error)
	// Thêm các hàm khác (GetLinkInfo, ListLinks...) sau này
}

// linkServiceImpl là triển khai cụ thể của LinkService
type linkServiceImpl struct {
	Repo repository.LinkRepository
	Cfg  *config.Config
}

func (l linkServiceImpl) CreateShortLink(originalURL string) (*model.LinkResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l linkServiceImpl) GetOriginalURL(shortCode string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// NewLinkService là hàm khởi tạo Tầng Service
func NewLinkService(repo repository.LinkRepository, cfg *config.Config) LinkService {
	return &linkServiceImpl{
		Repo: repo,
		Cfg:  cfg,
	}
}

// Định nghĩa lỗi nghiệp vụ
var ErrInvalidURL = errors.New("URL is invalid or empty")
var ErrLinkNotFound = errors.New("short link not found")
var ErrConflict = errors.New("short code conflict occurred, please retry")
var ErrServiceUnavailable = errors.New("internal service error")
