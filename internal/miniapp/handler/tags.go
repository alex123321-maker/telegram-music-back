package handler

import (
	"fmt"
	"strconv"
	"telegram-music/internal/miniapp/database"

	"github.com/gofiber/fiber/v3"
)

type TagRequest struct {
	Name string `json:"name" validate:"required"`
}

type TagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateTagHandler(c fiber.Ctx) error {
	var req TagRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный JSON: "+err.Error())
	}
	if req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Поле name обязательно")
	}

	tag, err := database.InsertTagIfNotExists(c.Context(), req.Name)
	if err != nil {
		return fiber.NewError(fiber.ErrBadGateway.Code, fmt.Sprintf("InsertTag: %v", err))
	}

	return c.JSON(TagResponse{ID: tag.ID, Name: tag.Name})
}

func ListTagsHandler(c fiber.Ctx) error {
	tags, err := database.GetAllTags(c.Context())
	if err != nil {
		return fiber.NewError(fiber.ErrBadGateway.Code, fmt.Sprintf("GetAllTags: %v", err))

	}

	resp := make([]TagResponse, len(tags))
	for i, t := range tags {
		resp[i] = TagResponse{ID: t.ID, Name: t.Name}
	}
	return c.JSON(resp)
}

func DeleteTagHandler(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный id")
	}

	if err := database.DeleteTag(c.Context(), id); err != nil {
		return fiber.NewError(fiber.ErrBadGateway.Code, fmt.Sprintf("DeleteTag: %v", err))

	}
	return c.SendStatus(fiber.StatusNoContent)
}
