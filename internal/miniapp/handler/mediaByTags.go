package handler

import (
	"fmt"

	"telegram-music/internal/miniapp/database"

	"github.com/gofiber/fiber/v3"
)

/* ---------------------- входной JSON ----------------------

{
  "tg_id":    123456789,      // обязательный int64
  "tags":     [1,2,3],        // опционально, []int
  "match_all": false          // опц., bool (по-умолчанию false)
}

------------------------------------------------------------*/

type MediaByTagsRequest struct {
	TgID     int64 `json:"tg_id"     validate:"required"`
	Tags     []int `json:"tags"`                // может отсутствовать
	MatchAll bool  `json:"match_all,omitempty"` // default = false
}

// POST /api/media
func GetMediaByTagsHandler(c fiber.Ctx) error {
	// ─── 1. читаем и валидируем тело ────────────────────────────────
	var req MediaByTagsRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "неверный JSON: "+err.Error())
	}
	if req.TgID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "`tg_id` должен быть > 0")
	}
	// (доп. проверка, что все tags > 0)
	for _, id := range req.Tags {
		if id <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "`tags` — только положительные числа")
		}
	}

	// ─── 2. бизнес-логика ───────────────────────────────────────────
	media, err := database.GetMediaByTagsExt(
		c.Context(),
		req.TgID,
		req.Tags,
		req.MatchAll,
	)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Sprintf("GetMediaByTagsExt: %v", err),
		)
	}

	return c.JSON(media)
}
