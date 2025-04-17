package miniapp

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	config "telegram-music/config/miniapp"
	"telegram-music/internal/miniapp/handler"
	"telegram-music/pkg/logging"
)

func Run(cfg *config.Config, logger logging.Logger) error {
	app := fiber.New()

	// 1) Логируем все входящие запросы
	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next() // передаём управление дальше
		logger.Infof(
			"%s %s → %d (%s)",
			c.Method(),
			c.OriginalURL(),
			c.Response().StatusCode(),
			time.Since(start),
		)
		return err
	})

	// 2) CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:5173", "https://mandrikov-ad.ru", "*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}))

	// 3) Маршруты
	registerRoutes(app)

	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Infof("Запуск miniapp на порту %d", cfg.Port)
	return app.Listen(addr)
}

func registerRoutes(app *fiber.App) {
	app.Get("/api/home", handler.HomeHandler)
	app.Post("/api/resolve", handler.ResolveMediaHandler)
}
