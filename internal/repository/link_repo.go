package repository

import (
	"errors"

	"github.com/quandang08/go-shorty/internal/model"
	"gorm.io/gorm"
)

// LinkRepository defines methods to interact with the Link table in the database.
type LinkRepository interface {
	Save(link *model.Link) error

	FindByShortCode(code string) (*model.Link, error)

	IncrementClicks(code string) error
}

// linkRepositoryImpl is the concrete implementation of LinkRepository.
type linkRepositoryImpl struct {
	DB *gorm.DB
}

// NewLinkRepository creates a new instance of LinkRepository.
func NewLinkRepository(db *gorm.DB) LinkRepository {
	return &linkRepositoryImpl{DB: db}
}

// Save inserts a new link record into the database.
func (l *linkRepositoryImpl) Save(link *model.Link) error {
	if err := l.DB.Create(link).Error; err != nil {
		return err
	}
	return nil
}

// FindByShortCode retrieves a link by its short code.
func (l *linkRepositoryImpl) FindByShortCode(code string) (*model.Link, error) {
	var link model.Link
	if err := l.DB.Where("short_code = ?", code).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &link, nil
}

// IncrementClicks increases the click count of a link.
func (l *linkRepositoryImpl) IncrementClicks(code string) error {
	result := l.DB.Model(&model.Link{}).
		Where("short_code = ?", code).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1))
	return result.Error
}
