BEGIN;

-- 1. Справочник тегов
CREATE TABLE tags (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- 2. Заменяем текстовый столбец на FK
ALTER TABLE media_tags
    ADD COLUMN tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    DROP COLUMN tag;

COMMIT;
