package links

import "time"

type Link struct {
	ID          int64
	OriginalURL string
	Code        string
	CreatedAt   time.Time
}

type StatsPeriod struct {
	From *time.Time
	To   *time.Time
}

type LinkStats struct {
	Code   string
	Clicks int64
}

type ClickEvent struct {
	LinkID    int64     `json:"link_id"`
	ClickedAt time.Time `json:"clicked_at"`
}

type ClickAggregate struct {
	LinkID int64
	Day    time.Time
	Clicks int64
}
