package handler

import (
	"fmt"

	"telegram-music/internal/miniapp/database"

	"github.com/gofiber/fiber/v3"
)

/* ---------------------- входной JSON ----------------------

{

  "tags":     [1,2,3],        // опционально, []int
  "match_all": false          // опц., bool (по-умолчанию false)
}

------------------------------------------------------------*/

type MediaByTagsRequest struct {
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

	for _, id := range req.Tags {
		if id <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "`tags` — только положительные числа")
		}
	}

	// ─── 2. бизнес-логика ───────────────────────────────────────────
	media, err := database.GetMediaByTagsExt(
		c.Context(),
		c.Locals("tg_id").(int64),
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
