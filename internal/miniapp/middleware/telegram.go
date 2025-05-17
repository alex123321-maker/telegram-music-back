package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

const (
	maxAge     = 5 * time.Minute // допустимый «срок давности»
	headerName = "X-Tg-Init-Data"
)

func CheckTelegram(botID int64) fiber.Handler {
	return func(c fiber.Ctx) error {
		var botID int64 = botID
		raw := c.Get(headerName)
		if raw == "" {
			return fiber.ErrUnauthorized
		}

		err := initdata.ValidateThirdParty(raw, botID, maxAge)
		if err != nil { // подпись не совпала
			return fiber.ErrUnauthorized
		}
		data, err := initdata.Parse(raw)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("tg_id", data.User.ID)
		return c.Next()
	}
}
