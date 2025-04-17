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

# базовые утилиты + python3 для yt-dlp
RUN apk add --no-cache \
      curl \
      ffmpeg \
      python3 \
      py3-cryptography \
  && adduser -D appuser

RUN curl -L \
  https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
-o /usr/local/bin/yt-dlp \
&& chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /app

# копируем бинарь сервера и миграции
COPY --from=builder /app/server .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/internal/miniapp/migrations ./migrations

# копируем entrypoint-скрипт
COPY deployments/miniapp/entrypoint.sh .
RUN chmod +x ./entrypoint.sh

USER appuser

EXPOSE 3000
ENTRYPOINT ["./entrypoint.sh"]