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
    
    # Установим зависимости (curl, yt-dlp и libc)
    RUN apk add --no-cache curl ffmpeg python3 py3-pip && \
        pip install yt-dlp && \
        adduser -D appuser
    
    WORKDIR /app
    
    # Копируем собранный бинарник из build-стадии
    COPY --from=builder /app/server .
    
    # Запуск от непривилегированного пользователя
    USER appuser
    
    # Порт, который будет использоваться внутри контейнера
    EXPOSE 3000
    
    # Запуск приложения
    CMD ["./server"]
    