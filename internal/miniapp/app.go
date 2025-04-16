package miniapp

import (
	"fmt"

	config "telegram-music/config/miniapp"
	"telegram-music/internal/miniapp/handler"
	"telegram-music/pkg/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Run(cfg *config.Config, logger logging.Logger) error {

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:5173, https://mandrikov-ad.ru"}, // Разрешённые фронтенды
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))
	registerRoutes(app)

	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Infof("Запуск miniapp на порту %d", cfg.Port)
	return app.Listen(addr)
}

func registerRoutes(app *fiber.App) {

	app.Get("api/home", handler.HomeHandler)
	app.Post("api/resolve", handler.ResolveMediaHandler)

}
