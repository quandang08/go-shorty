package model

import "time"

type Link struct {
	ID uint `gorm:"primaryKey" json:"-"`

	ShortCode string `gorm:"uniqueIndex;type:varchar(10)" json:"short_code"`

	OriginalURL string `gorm:"type:text;not null" json:"original_url"`

	ClicksCount uint `gorm:"default:0" json:"clicks_count"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
