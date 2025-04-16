CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    source_id TEXT NOT NULL,
    source_type TEXT NOT NULL,
    title TEXT NOT NULL,
    artist TEXT,
    description TEXT,
    url TEXT NOT NULL,             
    thumbnail_url TEXT,
    duration INT,                  
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (source_type, source_id)
);

CREATE TABLE media_tags (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT NOT NULL,
    media_id INTEGER REFERENCES media(id) ON DELETE CASCADE,
    tag TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
