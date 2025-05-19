package database

import (
	"context"
	"fmt"
	"time"
)

type Media struct {
	ID           int    // автоинкремент в БД
	SourceID     string // YouTube ID или аналог
	SourceType   string // например: "youtube"
	Title        string
	Artist       string
	Description  string
	URL          string
	ThumbnailURL string
	Duration     int
	OriginURL    string // новая ссылка на оригинал (origin_url)
	Exptime      int64  // время истечения в Unix
	CreatedAt    time.Time
}

func InsertOrUpdateMedia(ctx context.Context, m Media) (Media, error) {
	var id int
	var createdAt time.Time
	err := DB.QueryRow(ctx, `
		INSERT INTO media (
			source_id, source_type, title, artist, description, url,
			thumbnail_url, duration, origin_url, exptime
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (source_type, source_id) DO UPDATE SET
			title = EXCLUDED.title,
			artist = EXCLUDED.artist,
			description = EXCLUDED.description,
			url = EXCLUDED.url,
			thumbnail_url = EXCLUDED.thumbnail_url,
			duration = EXCLUDED.duration,
			origin_url = EXCLUDED.origin_url,
			exptime = EXCLUDED.exptime
		RETURNING id, created_at
	`,
		m.SourceID, m.SourceType, m.Title, m.Artist, m.Description, m.URL,
		m.ThumbnailURL, m.Duration, m.OriginURL, m.Exptime,
	).Scan(&id, &createdAt)

	m.ID = id
	m.CreatedAt = createdAt
	return m, err
}

type MediaWithTags struct {
	Media
	Tags []Tag `json:"tags"`
}

func GetMediaByTagsExt(
	ctx context.Context,
	tgID int64,
	tagIDs []int,
	matchAll bool,
) ([]MediaWithTags, error) {
	arr := make([]int32, len(tagIDs))
	for i, v := range tagIDs {
		arr[i] = int32(v)
	}

	cond := ""
	args := []any{tgID}

	switch {
	case len(tagIDs) == 0:
		// без условий
	case matchAll:
		cond = `
		  AND m.id IN (
		    SELECT  media_id
		    FROM    media_tags
		    WHERE   tg_id   = $1
		      AND   tag_id  = ANY($2)
		    GROUP BY media_id
		    HAVING   COUNT(DISTINCT tag_id) = $3
		  )`
		args = append(args, &arr, len(tagIDs))
	default:
		cond = `AND mt.tag_id = ANY($2)`
		args = append(args, &arr)
	}

	query := fmt.Sprintf(`
	  SELECT m.id, m.source_id, m.source_type,
	         m.title,
	         COALESCE(m.artist, '')        AS artist,
	         COALESCE(m.description, '')   AS description,
	         m.url,
	         COALESCE(m.thumbnail_url, '') AS thumbnail_url,
	         m.duration,
	         COALESCE(m.origin_url, '')    AS origin_url,
	         COALESCE(m.exptime, 0)        AS exptime,
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
		err := rows.Scan(
			&m.ID, &m.SourceID, &m.SourceType,
			&m.Title, &m.Artist, &m.Description,
			&m.URL, &m.ThumbnailURL,
			&m.Duration, &m.OriginURL, &m.Exptime,
			&m.CreatedAt,
			&tID, &tNm,
		)
		if err != nil {
			return nil, err
		}
		if currentMedia == nil || m.ID != currentID {
			currentMedia = &MediaWithTags{Media: m}
			res = append(res, *currentMedia)
			currentID = m.ID
		}
		tag := Tag{ID: int(tID), Name: tNm}
		currentMedia.Tags = append(currentMedia.Tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func GetMediaBySource(
	ctx context.Context,
	tgID int64, // идентификатор пользователя/чата
	sourceType string, // напр. "youtube"
	sourceID string, // YouTube-ID
) (MediaWithTags, error) {

	const q = `
	  SELECT m.id, m.source_id, m.source_type,
	         m.title,
	         COALESCE(m.artist,'')        AS artist,
	         COALESCE(m.description,'')   AS description,
	         m.url,
	         COALESCE(m.thumbnail_url,'') AS thumbnail_url,
	         m.duration,
	         COALESCE(m.origin_url,'')    AS origin_url,
	         COALESCE(m.exptime,0)        AS exptime,
	         m.created_at,
	         COALESCE(t.id, 0)            AS tag_id,   -- может быть NULL, но Scan требует int
	         COALESCE(t.name,'')          AS tag_name
	  FROM   media m
	  LEFT JOIN media_tags mt ON mt.media_id = m.id
	                           AND mt.tg_id   = $3
	  LEFT JOIN tags t        ON t.id        = mt.tag_id
	  WHERE  m.source_type = $1
	    AND  m.source_id   = $2;
	`

	rows, err := DB.Query(ctx, q, sourceType, sourceID, tgID)
	if err != nil {
		return MediaWithTags{}, err
	}
	defer rows.Close()

	var (
		result   MediaWithTags
		initDone bool
	)

	for rows.Next() {
		var (
			m   Media
			tID int
			tNm string
		)
		if err := rows.Scan(
			&m.ID, &m.SourceID, &m.SourceType,
			&m.Title, &m.Artist, &m.Description,
			&m.URL, &m.ThumbnailURL,
			&m.Duration, &m.OriginURL, &m.Exptime,
			&m.CreatedAt,
			&tID, &tNm,
		); err != nil {
			return MediaWithTags{}, err
		}

		if !initDone {
			result.Media = m
			initDone = true
		}
		// если тегов нет, tID будет 0 и tNm пустой ─ пропускаем
		if tID != 0 {
			result.Tags = append(result.Tags, Tag{ID: tID, Name: tNm})
		}
	}
	if err = rows.Err(); err != nil {
		return MediaWithTags{}, err
	}
	if !initDone {
		return MediaWithTags{}, fmt.Errorf("media с source_type=%s, source_id=%s не найдено", sourceType, sourceID)
	}
	return result, nil
}
