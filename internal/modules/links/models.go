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
