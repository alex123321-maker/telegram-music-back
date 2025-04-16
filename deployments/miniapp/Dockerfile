# ==================================================================
#                          Build Stage
# ==================================================================
FROM golang:1.24-alpine AS builder

# установим git для go mod и миграций
RUN apk add --no-cache git

WORKDIR /app

# копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# копируем весь исходник
COPY . .

# собираем ваше приложение
RUN go build -o server ./cmd/miniapp

# ставим golang-migrate CLI нужной версии
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2


# ==================================================================
#                         Runtime Stage
# ==================================================================
FROM alpine:3.19

# ставим только то, что нужно для работы: yt-dlp, ffmpeg, curl
RUN apk add --no-cache curl ffmpeg yt-dlp \
    && adduser -D appuser

WORKDIR /app

# копируем ваш собранный сервер
COPY --from=builder /app/server .

# копируем утилиту migrate, установленную во время сборки
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# копируем папку с миграциями
COPY --from=builder /app/internal/miniapp/migrations ./migrations

# снижаем привилегии
USER appuser

EXPOSE 3000

# При старте сначала применяем все миграции, затем запускаем сервер
ENTRYPOINT ["sh", "-c", \
  "migrate -path ./migrations -database \"$DATABASE_URL\" up && exec ./server"]
