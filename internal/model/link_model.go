package model

import "time"

type Link struct {
	ID          int64     `db:"id" json:"id"`
	ShortURL    string    `db:"short_url" json:"short_url"`
	OriginalURL string    `db:"original_url" json:"original_url"`
	Clicks      int64     `db:"clicks" json:"clicks"`
	OwnerID     int64     `db:"owner_id" json:"owner_id"`
	ExpiredAt   time.Time `db:"expired_at" json:"expired_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
