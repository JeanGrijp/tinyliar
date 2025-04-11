package repository

import (
	"github.com/JeanGrijp/tinyliar/internal/model"
	"github.com/jmoiron/sqlx"
)

type LinkRepository struct {
	db *sqlx.DB
}

func NewLinkRepository(db *sqlx.DB) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (r *LinkRepository) CreateLink(link *model.Link) error {
	query := `INSERT INTO links (short_url, original_url, clicks, owner_id, expired_at, created_at, updated_at) 
	VALUES (:short_url, :original_url, :clicks, :owner_id, :expired_at, :created_at, :updated_at)`
	_, err := r.db.NamedExec(query, link)
	if err != nil {
		return err
	}
	return nil
}

func (r *LinkRepository) GetLinkByShortURL(shortURL string) (*model.Link, error) {
	query := `SELECT * FROM links WHERE short_url = ?`
	link := &model.Link{}
	err := r.db.Get(link, query, shortURL)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *LinkRepository) IncrementClickCount(id int64) error {
	query := `UPDATE links SET clicks = clicks + 1 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
