package repository

import (
	"errors"

	"github.com/quandang08/go-shorty/internal/model"
	"gorm.io/gorm"
)

// LinkRepository defines the DB operations for Link.
type LinkRepository interface {
	Create(link *model.Link) error
	UpdateShortCode(link *model.Link) error
	FindByShortCode(code string) (*model.Link, error)
	IncrementClicks(code string) error
	FindAll() ([]model.Link, error)
}

type linkRepositoryImpl struct {
	DB *gorm.DB
}

func NewLinkRepository(db *gorm.DB) LinkRepository {
	return &linkRepositoryImpl{DB: db}
}

// Create inserts a new Link record and populates its ID.
func (r *linkRepositoryImpl) Create(link *model.Link) error {
	return r.DB.Create(link).Error
}

// UpdateShortCode updates only the short_code column.
func (r *linkRepositoryImpl) UpdateShortCode(link *model.Link) error {
	result := r.DB.Model(&model.Link{}).
		Where("id = ?", link.ID).
		UpdateColumn("short_code", link.ShortCode)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

// FindByShortCode retrieves a Link record by short code.
func (r *linkRepositoryImpl) FindByShortCode(code string) (*model.Link, error) {
	var link model.Link
	err := r.DB.Where("short_code = ?", code).First(&link).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &link, nil
}

// IncrementClicks atomically increases the click counter.
func (r *linkRepositoryImpl) IncrementClicks(code string) error {
	return r.DB.Model(&model.Link{}).
		Where("short_code = ?", code).
		UpdateColumn("clicks_count", gorm.Expr("clicks_count + ?", 1)).
		Error
}

// FindAll retrieves all link records.
func (r *linkRepositoryImpl) FindAll() ([]model.Link, error) {
	var links []model.Link
	if result := r.DB.Find(&links); result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}
