#!/bin/sh
set -e

# 1) применяем миграции
migrate -path ./migrations -database "$DATABASE_URL" up

# 2) запускаем сервер
exec ./server
