#!/bin/sh
set -e

# 1) скачиваем самую свежую сборку yt-dlp из GitHub Releases
curl -L \
  https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
  -o /usr/local/bin/yt-dlp \
  && chmod a+rx /usr/local/bin/yt-dlp

# 2) применяем миграции
migrate -path ./migrations -database "$DATABASE_URL" up

# 3) запускаем сервер
exec ./server
