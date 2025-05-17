package handler

import (
	"fmt"
	"strconv"
	"time"

	"telegram-music/internal/miniapp/database"

	"github.com/gofiber/fiber/v3"
)

/*--------------------------------------------------------------
   DTO
--------------------------------------------------------------*/

type LinkTagRequest struct {
	MediaID int `json:"media_id"` // обязателен
	TagID   int `json:"tag_id"`   // обязателен
}

type MediaTagResponse struct {
	ID        int       `json:"id"`
	TgID      int64     `json:"tg_id"`
	MediaID   int       `json:"media_id"`
	TagID     int       `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

/*
--------------------------------------------------------------

	POST /api/media-tags  — добавить связь

--------------------------------------------------------------
*/
func CreateMediaTagHandler(c fiber.Ctx) error {
	var req LinkTagRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "неверный JSON: "+err.Error())
	}

	// валидация
	if req.MediaID <= 0 || req.TagID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "media_id и tag_id должны быть > 0")
	}

	mt, err := database.InsertMediaTagIfNotExists(
		c.Context(), c.Locals("tg_id").(int64), req.MediaID, req.TagID,
	)
	if err != nil {
		return fiber.NewError(fiber.ErrBadGateway.Code, fmt.Sprintf("Не удалось добавить тег: %s", err))

	}

	return c.Status(fiber.StatusCreated).JSON(MediaTagResponse{
		ID:        mt.ID,
		TgID:      mt.TgID,
		MediaID:   mt.MediaID,
		TagID:     mt.TagID,
		CreatedAt: mt.CreatedAt,
	})
}

/*
--------------------------------------------------------------

	DELETE /api/media-tags/:id  — удалить связь

--------------------------------------------------------------
*/
func DeleteMediaTagHandler(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "некорректный id")
	}

	if err := database.DeleteMediaTag(c.Context(), id); err != nil {
		return fiber.NewError(fiber.ErrBadGateway.Code, fmt.Sprintf("Не удалось убрать тег: %s", err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetUserTagsHandler(c fiber.Ctx) error {
	tags, err := database.GetTagsForUser(c.Context(), c.Locals("tg_id").(int64))
	if err != nil {
		return fiber.NewError(fiber.ErrBadGateway.Code, fmt.Sprintf("Ошибка получения тегов: %s", err))
	}
	return c.JSON(tags)
}
