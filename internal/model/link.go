package model

import "time"

// Link represents a shortened URL entity stored in the database.
type Link struct {
	ID          uint      `gorm:"primaryKey" json:"-"`
	ShortCode   string    `gorm:"uniqueIndex;type:varchar(10)" json:"short_code"`
	OriginalURL string    `gorm:"type:text;not null" json:"original_url"`
	ClicksCount uint      `gorm:"default:0" json:"clicks_count"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// CreateLinkRequest is the input payload sent by the client when requesting
// a new shortened link.
type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" form:"original_url" binding:"required"`
}

// LinkResponse represents the response payload returned to the client
// after a short link has been created.
type LinkResponse struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ClicksCount uint      `json:"clicks_count"`
	CreatedAt   time.Time `json:"created_at"`
	ShortURL    string    `json:"short_url"`
}
