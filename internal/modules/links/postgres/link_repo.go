package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *LinkRepo) RecordClick(ctx context.Context, linkID int64, clickedAt time.Time) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO link_clicks (link_id, clicked_at) VALUES ($1, $2)`,
		linkID,
		clickedAt,
	)
	if err != nil {
		return fmt.Errorf("insert link click: %w", err)
	}
	return nil
}

func (r *LinkRepo) GetStats(ctx context.Context, code string, period links.StatsPeriod) (*links.LinkStats, error) {
	stats := &links.LinkStats{}
	err := r.db.QueryRow(
		ctx,
		`SELECT l.code, COUNT(c.id)
		 FROM links l
		 LEFT JOIN link_clicks c ON c.link_id = l.id
		   AND ($2::timestamptz IS NULL OR c.clicked_at >= $2)
		   AND ($3::timestamptz IS NULL OR c.clicked_at <= $3)
		 WHERE l.code = $1
		 GROUP BY l.code`,
		code,
		period.From,
		period.To,
	).Scan(&stats.Code, &stats.Clicks)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, links.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query link stats: %w", err)
	}
	return stats, nil
}
