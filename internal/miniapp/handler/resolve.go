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
		return fiber.NewError(fiber.StatusBadRequest, "Неверный JSON: "+err.Error())
	}

	// 1) Пытаемся извлечь YouTube-ID
	sourceID, err := service.ExtractYouTubeID(req.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	const sourceType = "youtube"

	ctx := context.Background()

	// 2) Проверяем, есть ли запись в БД
	mediaDB, err := database.GetMediaBySource(ctx, 0, sourceType, sourceID) // tgID = 0 (теги не нужны)
	foundInDB := err == nil

	// 3) Если запись найдена и ссылка ещё жива — сразу отдаём
	if foundInDB && !service.IsExpiredSoon(mediaDB.Exptime, mediaDB.Duration) {
		return c.JSON(mediaDB)
	}

	resolved, err := service.ResolveUrl(req.URL)
	if err != nil {
		return fiber.NewError(fiber.ErrBadGateway.Code, err.Error())
	}

	// 5) Пишем в БД (InsertOrUpdate)
	media := database.Media{
		SourceID:     resolved.SourceID,
		SourceType:   resolved.SourceType,
		Title:        resolved.Title,
		Artist:       resolved.Artist,
		Description:  resolved.Description,
		URL:          resolved.AudioURL,
		ThumbnailURL: resolved.ThumbnailURL,
		Duration:     resolved.Duration,
		OriginURL:    resolved.OriginURL,
		Exptime:      resolved.Exptime,
	}

	media, err = database.InsertOrUpdateMedia(ctx, media)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка сохранения в БД: "+err.Error())
	}

	// 6) Отдаём обновлённый объект
	return c.JSON(media)
}
