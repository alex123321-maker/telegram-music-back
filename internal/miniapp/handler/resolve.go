package handler

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	database "telegram-music/internal/miniapp/database"

	"github.com/gofiber/fiber/v3"
)

type ResolveMediaRequest struct {
	URL string `json:"url"`
}

type ResolveMediaResponse struct {
	ID           int     `json:"id"`
	SourceID     string  `json:"source_id"`
	SourceType   string  `json:"source_type"`
	Title        string  `json:"title"`
	Artist       *string `json:"artist,omitempty"`
	Description  *string `json:"description,omitempty"`
	Duration     int     `json:"duration"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
	AudioURL     string  `json:"audio_url"` // ← прямая ссылка на аудио
}

func ResolveMediaHandler(c fiber.Ctx) error {
	var req ResolveMediaRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный JSON")
	}

	ctx := context.Background()

	cmd := exec.CommandContext(ctx, "yt-dlp",
		"--quiet", "--no-warnings",
		"--no-playlist",
		"--format", "bestaudio[ext=m4a]/bestaudio",
		"--print", "%(id)s␞%(title)s␞%(duration)s␞%(uploader)s␞%(description)s␞%(thumbnail)s␞%(webpage_url)s␞%(url)s",
		req.URL,
	)

	output, err := cmd.Output()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "yt-dlp не смог обработать ссылку: "+err.Error())
	}
	fmt.Println("YT-DLP raw output:", string(output))
	separator := "␞" /
	parts := strings.Split(string(output), separator)
	if len(parts) < 8 {
		return fiber.NewError(fiber.StatusInternalServerError, "yt-dlp вернул недостаточно данных")
	}

	duration, _ := strconv.Atoi(parts[2])

	// nullable поля
	var (
		artistPtr      *string
		descriptionPtr *string
		thumbnailPtr   *string
	)

	if parts[3] != "" {
		artist := strings.TrimSpace(parts[3])
		artistPtr = &artist
	}
	if parts[4] != "" {
		desc := strings.TrimSpace(parts[4])
		descriptionPtr = &desc
	}
	if parts[5] != "" {
		thumb := strings.TrimSpace(parts[5])
		thumbnailPtr = &thumb
	}

	audioURL := strings.TrimSpace(parts[7])
	if audioURL == "" || audioURL == "NA" {
		return fiber.NewError(fiber.StatusInternalServerError, "yt-dlp не смог извлечь аудиоссылку")
	}

	media := database.Media{
		SourceID:     parts[0],
		SourceType:   "youtube",
		Title:        parts[1],
		Artist:       artistPtr,
		Description:  descriptionPtr,
		URL:          parts[6],
		ThumbnailURL: thumbnailPtr,
		Duration:     duration,
		CreatedAt:    time.Now(),
	}
	// вставка в базу
	id, err := database.InsertMediaIfNotExists(ctx, media)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при сохранении в базу: "+err.Error())
	}
	return c.JSON(ResolveMediaResponse{
		ID:           id, // можно заменить на media.ID после сохранения
		SourceID:     media.SourceID,
		SourceType:   media.SourceType,
		Title:        media.Title,
		Artist:       media.Artist,
		Description:  media.Description,
		Duration:     media.Duration,
		ThumbnailURL: media.ThumbnailURL,
		AudioURL:     audioURL,
	})
}
