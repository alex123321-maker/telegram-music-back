package main

import (
	"context"
	"log"
	config "telegram-music/config/miniapp"
	miniapp "telegram-music/internal/miniapp"
	database "telegram-music/internal/miniapp/database"
	"telegram-music/pkg/logging"
)

func main() {
	// Загружаем конфиг (например, через переменные окружения)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	// Инициализируем логгер
	logger := logging.NewLogger(cfg)

	ctx := context.Background()
	if err := database.InitDB(ctx, cfg.DatabaseURL); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer database.DB.Close()
	// Запускаем само приложение (Fiber)
	if err := miniapp.Run(cfg, logger); err != nil {
		logger.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
