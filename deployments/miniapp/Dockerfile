# ---------- Build Stage ----------
    FROM golang:1.24-alpine AS builder

    # Установим зависимости
    RUN apk add --no-cache git
    
    # Создаём рабочую директорию
    WORKDIR /app
    
    # Копируем и устанавливаем зависимости
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Копируем исходники
    COPY . .
    
    # Сборка
    RUN go build -o server ./cmd/miniapp
    
    
    # ---------- Runtime Stage ----------
    FROM alpine:3.21

    RUN apk add --no-cache curl ffmpeg yt-dlp \
        && apk add --no-cache --virtual .build-deps curl git \
        && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz \
        | tar xz -C /usr/local/bin migrate.linux-amd64 \
        && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate \
        && apk del .build-deps

    WORKDIR /app
    COPY --from=builder /app/server .
    COPY --from=builder /app/internal/miniapp/migrations ./migrations

    USER appuser
    EXPOSE 3000

    # при старте сначала миграции, затем ваш сервер
    ENTRYPOINT ["sh", "-c", "migrate -path ./migrations -database \"$DATABASE_URL\" up && exec ./server"]


    