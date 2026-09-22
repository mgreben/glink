package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mgreben/glink/internal/modules/links"
)

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinkRepo(db *pgxpool.Pool) *LinkRepo {
	return &LinkRepo{
		db: db,
	}
}

func (r *LinkRepo) GetByCode(ctx context.Context, code string) (*links.Link, error) {
	link := &links.Link{}

	err := r.db.QueryRow(
		ctx,
		`SELECT id, original_url, code, created_at
		 FROM links
		 WHERE code = $1`,
		code,
	).Scan(&link.ID, &link.OriginalURL, &link.Code, &link.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, links.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (r *LinkRepo) Create(ctx context.Context, link *links.Link) error {
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO links (original_url, code, created_at)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		link.OriginalURL,
		link.Code,
		link.CreatedAt,
	).Scan(&link.ID)

	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return links.ErrUniqueCode
	}

	return err
}
