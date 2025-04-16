package handler

import (
	"github.com/gofiber/fiber/v3"
)

func HomeHandler(c fiber.Ctx) error {
	return c.SendString("Это главная страница мини-приложения.")
}
