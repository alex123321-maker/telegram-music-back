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
	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next() // даём обработать запрос/ошибку

		status := c.Response().StatusCode()
		if err != nil { // если вернулась ошибка —
			if fe, ok := err.(*fiber.Error); ok {
				status = fe.Code
			} else {
				status = fiber.StatusInternalServerError
			}
		}

		logger.Infof("%s %s → %d (%s)",
			c.Method(), c.OriginalURL(), status, time.Since(start))

		return err
	})
	// app.Use(middleware.CheckTelegram(cfg.BotId))
	app.Use(func(c fiber.Ctx) error {
		c.Locals("tg_id", int64(1893384316)) // Замените на ваш реальный tg_id
		return c.Next()
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:5173", "https://mandrikov-ad.ru", "*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}))

	registerRoutes(app)

	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Infof("Запуск miniapp на порту %d", cfg.Port)
	return app.Listen(addr)
}

func registerRoutes(app *fiber.App) {
	app.Post("/api/resolve", handler.ResolveMediaHandler)
	app.Post("/api/media", handler.GetMediaByTagsHandler)
	app.Post("/api/tags", handler.CreateTagHandler)
	app.Get("/api/tags", handler.ListTagsHandler)
	app.Get("api/media/:media_id/my-tags", handler.GetMediaTagsHandler)
	app.Get("/api/my-tags", handler.GetUserTagsHandler)
	app.Delete("/api/tags/:id", handler.DeleteTagHandler)
}
