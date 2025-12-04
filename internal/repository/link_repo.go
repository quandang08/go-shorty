package repository

import (
	"github.com/quandang08/go-shorty/internal/model"
	"gorm.io/gorm"
)

// LinkRepository định nghĩa các phương thức giao tiếp với DB

type LinkRepository interface {
	Save(link *model.Link) error
	FindByShortCode(code string) (*model.Link, error)
	IncrementClicks(code string) error

	// Thêm các hàm khác (FindByLinkID, FindAll...) sau này.
}

type linkRepositoryImpl struct {
	DB *gorm.DB
}

func (l linkRepositoryImpl) Save(link *model.Link) error {
	//TODO implement me
	panic("implement me")
}

func (l linkRepositoryImpl) FindByShortCode(code string) (*model.Link, error) {
	//TODO implement me
	panic("implement me")
}

func (l linkRepositoryImpl) IncrementClicks(code string) error {
	//TODO implement me
	panic("implement me")
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
	return &linkRepositoryImpl{DB: db}
}

// --- Các hàm triển khai (Save, FindByShortCode, IncrementClicks) sẽ nằm dưới đây ---
