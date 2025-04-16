package miniapp

import (
	"fmt"

	config "telegram-music/config/miniapp"
	"telegram-music/internal/miniapp/handler"
	"telegram-music/pkg/logging"

	"github.com/gofiber/fiber/v3"
)

func Run(cfg *config.Config, logger logging.Logger) error {

	app := fiber.New()

	registerRoutes(app)

	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Infof("Запуск miniapp на порту %d", cfg.Port)
	return app.Listen(addr)
}

func registerRoutes(app *fiber.App) {

	app.Get("api/home", handler.HomeHandler)
	app.Post("api/resolve", handler.ResolveMediaHandler)

}
