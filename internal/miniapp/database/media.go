package database

import (
	"context"
	"time"
)

type Media struct {
	ID           int    // автоинкремент в БД, можно оставить 0 при вставке
	SourceID     string // YouTube ID или аналог из другого источника
	SourceType   string // например: "youtube"
	Title        string
	Artist       *string // uploader, nullable
	Description  *string // optional
	URL          string  // прямая ссылка на оригинал (webpage_url)
	ThumbnailURL *string // прямой thumbnail (если есть)
	Duration     int     // в секундах
	CreatedAt    time.Time
}

func InsertMediaIfNotExists(ctx context.Context, m Media) (int, error) {
	var id int
	err := DB.QueryRow(ctx, `
		INSERT INTO media (
			source_id, source_type, title, artist, description, url, thumbnail_url, duration
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (source_type, source_id) DO UPDATE SET
			title = EXCLUDED.title
		RETURNING id
	`, m.SourceID, m.SourceType, m.Title, m.Artist, m.Description, m.URL, m.ThumbnailURL, m.Duration).Scan(&id)
	return id, err
}
