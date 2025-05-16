BEGIN;

-- 1. Возвращаем старый столбец
ALTER TABLE media_tags
    ADD COLUMN tag TEXT;

-- (таблица была пустая, поэтому заполнять нечего)

-- 2. Удаляем FK-столбец
ALTER TABLE media_tags
    DROP COLUMN tag_id;

-- 3. Удаляем справочник тегов
DROP TABLE tags;

COMMIT;
