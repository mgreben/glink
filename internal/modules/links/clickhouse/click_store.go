package clickhouse

import (
	"context"
	"fmt"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/mgreben/glink/internal/modules/links"
)

type ClickStore struct{ db ch.Conn }

func NewClickStore(db ch.Conn) *ClickStore { return &ClickStore{db: db} }

func (s *ClickStore) InsertClickEvents(ctx context.Context, events []links.ClickEvent) error {
	if len(events) == 0 {
		return nil
	}
	batch, err := s.db.PrepareBatch(ctx, "INSERT INTO link_clicks (link_id, clicked_at)")
	if err != nil {
		return fmt.Errorf("prepare click batch: %w", err)
	}
	for _, event := range events {
		if err := batch.Append(event.LinkID, event.ClickedAt); err != nil {
			return err
		}
	}
	if err := batch.Send(); err != nil {
		return fmt.Errorf("insert click batch: %w", err)
	}
	return nil
}

func (s *ClickStore) CountClicks(ctx context.Context, linkID int64, period links.StatsPeriod) (int64, error) {
	var clicks int64
	err := s.db.QueryRow(ctx, "SELECT count() FROM link_clicks WHERE link_id = ? AND (? IS NULL OR clicked_at >= ?) AND (? IS NULL OR clicked_at <= ?)", linkID, period.From, period.From, period.To, period.To).Scan(&clicks)
	return clicks, err
}
