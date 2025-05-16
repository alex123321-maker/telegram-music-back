package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Tag struct {
	ID   int
	Name string
}

// -------------------------------------------------------------------
// CRUD-ФУНКЦИИ
// -------------------------------------------------------------------

// InsertTagIfNotExists создаёт тег (name UNIQUE) или возвращает существующий.
func InsertTagIfNotExists(ctx context.Context, name string) (Tag, error) {
	const q = `
		INSERT INTO tags (name)
		VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id;
	`

	var id int
	if err := DB.QueryRow(ctx, q, name).Scan(&id); err != nil {
		return Tag{}, err
	}
	return Tag{ID: id, Name: name}, nil
}

// GetAllTags возвращает все теги по алфавиту.
func GetAllTags(ctx context.Context) ([]Tag, error) {
	rows, err := DB.Query(ctx, `SELECT id, name FROM tags ORDER BY name;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

// ErrTagNotFound используется при обращении к несуществующему тегу.
var ErrTagNotFound = errors.New("Тэг не найден")

// GetTagByID возвращает тег по ID или ErrTagNotFound.
func GetTagByID(ctx context.Context, id int) (Tag, error) {
	var t Tag
	err := DB.QueryRow(ctx, `SELECT id, name FROM tags WHERE id = $1;`, id).
		Scan(&t.ID, &t.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return Tag{}, ErrTagNotFound
	}
	return t, err
}

// DeleteTag удаляет тег по ID (привязки в media_tags исчезнут благодаря ON DELETE CASCADE).
func DeleteTag(ctx context.Context, id int) error {
	_, err := DB.Exec(ctx, `DELETE FROM tags WHERE id = $1;`, id)
	return err
}
