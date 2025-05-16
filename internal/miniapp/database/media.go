package database

import (
	"context"
	"fmt"
	"time"
)

type Media struct {
	ID           int    // автоинкремент в БД, можно оставить 0 при вставке
	SourceID     string // YouTube ID или аналог из другого источника
	SourceType   string // например: "youtube"
	Title        string
	Artist       string // uploader, nullable
	Description  string // optional
	URL          string // прямая ссылка на оригинал (webpage_url)
	ThumbnailURL string // прямой thumbnail (если есть)
	Duration     int    // в секундах
	CreatedAt    time.Time
}

func InsertMediaIfNotExists(ctx context.Context, m Media) (Media, error) {
	var id int
	var created_at time.Time
	err := DB.QueryRow(ctx, `
		INSERT INTO media (
			source_id, source_type, title, artist, description, url, thumbnail_url, duration
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (source_type, source_id) DO UPDATE SET
			title = EXCLUDED.title
		RETURNING id,created_at
	`, m.SourceID, m.SourceType, m.Title, m.Artist, m.Description, m.URL, m.ThumbnailURL, m.Duration).Scan(&id, &created_at)
	m.CreatedAt = created_at
	m.ID = id
	return m, err
}

type MediaWithTags struct {
	Media       // база
	Tags  []Tag `json:"tags"` // список тег-структур
}

// -------------------------------------------------------------------
// 2. Выборка
// -------------------------------------------------------------------

func GetMediaByTagsExt(
	ctx context.Context,
	tgID int64,
	tagIDs []int, // набор тегов фильтра
	matchAll bool, // true → нужны все теги; false → любой
) ([]MediaWithTags, error) {

	/* ---------------------------------------------------------------
	   Стратегия:
	   ─ Формируем один SELECT, который отдаёт
	     media-поля + (id, name) текущего тега.
	   ─ В Go-коде агрегируем строки по media.id и
	     собираем срез []Tag для каждого медиа.
	----------------------------------------------------------------- */

	// подготовим массив тегов для ANY() / = ALL()
	arr := make([]int32, len(tagIDs))
	for i, v := range tagIDs {
		arr[i] = int32(v)
	}

	cond := ``
	args := []any{tgID} // $1 всегда tg_id

	switch {
	case len(tagIDs) == 0:
		// без фильтра по тегам
	case matchAll:
		// медиa, содержащие ВСЕ указанные теги
		cond = `
		  AND m.id IN (
		    SELECT  media_id
		    FROM    media_tags
		    WHERE   tg_id   = $1
		      AND   tag_id  = ANY($2)
		    GROUP BY media_id
		    HAVING   COUNT(DISTINCT tag_id) = $3
		  )`
		args = append(args, &arr, len(tagIDs)) // $2, $3
	default:
		// медиa, содержащие хотя бы ОДИН тег
		cond = `AND mt.tag_id = ANY($2)`
		args = append(args, &arr) // $2
	}

	query := fmt.Sprintf(`
	  SELECT m.id, m.source_id, m.source_type,
	         m.title,
	         COALESCE(m.artist,'')       AS artist,
	         COALESCE(m.description,'')  AS description,
	         m.url,
	         COALESCE(m.thumbnail_url,'') AS thumbnail_url,
	         m.duration,
	         m.created_at,
	         t.id  AS tag_id,
	         t.name AS tag_name
	  FROM   media       m
	  JOIN   media_tags  mt ON mt.media_id = m.id
	  JOIN   tags        t  ON t.id        = mt.tag_id
	  WHERE  mt.tg_id = $1
	  %s
	  ORDER BY m.id, t.name;
	`, cond)

	rows, err := DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// ---------- агрегируем в Go ----------
	var (
		res          []MediaWithTags
		currentID    int
		currentMedia *MediaWithTags
	)

	for rows.Next() {
		var (
			m   Media
			tID int32
			tNm string
		)
		if err := rows.Scan(
			&m.ID, &m.SourceID, &m.SourceType,
			&m.Title, &m.Artist, &m.Description,
			&m.URL, &m.ThumbnailURL,
			&m.Duration, &m.CreatedAt,
			&tID, &tNm,
		); err != nil {
			return nil, err
		}

		// новая медиазапись?
		if currentMedia == nil || m.ID != currentID {
			currentMedia = &MediaWithTags{Media: m}
			res = append(res, *currentMedia)
			currentID = m.ID
		}
		// добавляем тег к текущему медиа
		tag := Tag{ID: int(tID), Name: tNm}
		currentMedia.Tags = append(currentMedia.Tags, tag)
	}

	// rows.Err()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
