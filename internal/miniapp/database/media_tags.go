package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type MediaTag struct {
	ID        int
	TgID      int64
	MediaID   int
	TagID     int
	CreatedAt time.Time
}

// InsertMediaTagIfNotExists добавляет связь
// «телеграм-юзер ↔ медиа ↔ тег».
// Если такая пара уже есть,
// возвращает существующую строку.
func InsertMediaTagIfNotExists(
	ctx context.Context,
	tgID int64,
	mediaID, tagID int,
) (MediaTag, error) {

	const q = `
		INSERT INTO media_tags (tg_id, media_id, tag_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (tg_id, media_id, tag_id) DO UPDATE
		SET tg_id = EXCLUDED.tg_id   
		RETURNING id, created_at;
	`

	var mt MediaTag
	mt.TgID, mt.MediaID, mt.TagID = tgID, mediaID, tagID

	err := DB.QueryRow(ctx, q, tgID, mediaID, tagID).
		Scan(&mt.ID, &mt.CreatedAt)

	// если возврата нет (старые версии PG),
	// вытягиваем существующую строку
	if errors.Is(err, pgx.ErrNoRows) {
		err = DB.QueryRow(
			ctx,
			`SELECT id, created_at FROM media_tags
			 WHERE tg_id=$1 AND media_id=$2 AND tag_id=$3`,
			tgID, mediaID, tagID,
		).Scan(&mt.ID, &mt.CreatedAt)
	}
	return mt, err
}

type MediaTagWithName struct {
	TagID int
	Name  string
}

func GetTagsForMedia(ctx context.Context, mediaID int) ([]MediaTagWithName, error) {
	const q = `
		SELECT mt.tag_id, t.name
		FROM media_tags mt
		JOIN tags t ON t.id = mt.tag_id
		WHERE mt.media_id = $1
		ORDER BY t.name;
	`

	rows, err := DB.Query(ctx, q, mediaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []MediaTagWithName
	for rows.Next() {
		var r MediaTagWithName
		if err := rows.Scan(&r.TagID, &r.Name); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, rows.Err()
}

// DeleteMediaTag удаляет строку
// по её первичному ключу id.
func DeleteMediaTag(ctx context.Context, id int) error {
	_, err := DB.Exec(ctx, `DELETE FROM media_tags WHERE id = $1;`, id)
	return err
}
