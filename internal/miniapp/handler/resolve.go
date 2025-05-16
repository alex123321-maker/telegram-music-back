package handler

import (
	"context"
	"telegram-music/internal/miniapp/database"
	"telegram-music/internal/miniapp/service"

	"github.com/gofiber/fiber/v3"
)

type ResolveMediaRequest struct {
	URL string `json:"url"`
}

func ResolveMediaHandler(c fiber.Ctx) error {
	var req ResolveMediaRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный JSON:"+err.Error())
	}
	resp, err := service.ResolveUrl(req.URL)
	if err != nil {
		return fiber.NewError(fiber.ErrBadGateway.Code, err.Error())
	}

	// 5) Сохраняем в базу
	media := database.Media{
		SourceID:     resp.SourceID,
		SourceType:   resp.SourceType,
		Title:        resp.Title,
		Artist:       resp.Artist,
		Description:  resp.Description,
		URL:          resp.AudioURL,
		ThumbnailURL: resp.ThumbnailURL,
		Duration:     resp.Duration,
	}

	media, err = database.InsertMediaIfNotExists(context.Background(), media)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError,
			"Ошибка сохранения в БД: "+err.Error())
	}

	return c.JSON(media)
}
