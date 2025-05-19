package service

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	ytdlp "github.com/lrstanley/go-ytdlp"
)

type ResolveMediaResponse struct {
	ID           int    `json:"id"`
	SourceID     string `json:"source_id"`
	SourceType   string `json:"source_type"`
	Title        string `json:"title"`
	Artist       string `json:"artist"`
	Description  string `json:"description"`
	Duration     int    `json:"duration"`
	ThumbnailURL string `json:"thumbnail_url"`
	OriginURL    string `json:"origin_url"`
	AudioURL     string `json:"audio_url"`
	Exptime      int64  `json:"exptime"`
}

func ResolveUrl(req string) (ResolveMediaResponse, error) {
	dl := ytdlp.New().
		Quiet().
		NoWarnings().
		NoPlaylist().
		Format("bestaudio[ext=m4a]/bestaudio").
		Print("%(id)s␞%(title)s␞%(duration)s␞%(uploader)s␞%(description)s␞%(thumbnail)s␞%(webpage_url)s␞%(url)s")
	raw, err := dl.Run(context.TODO(), req)
	if err != nil {
		return ResolveMediaResponse{}, fmt.Errorf("ошибка при полчении ответа от dlp")
	}

	parts := strings.Split(string(raw.Stdout), "␞")
	if len(parts) != 8 {
		return ResolveMediaResponse{}, fmt.Errorf("yt-dlp вернул неверный формат, %d", len(parts))
	}

	dur, _ := strconv.Atoi(parts[2])

	exptime, err := parseExptime(parts[7])
	if err != nil {
		return ResolveMediaResponse{}, fmt.Errorf("Ошибка:%s", err)
	}

	resp := ResolveMediaResponse{
		SourceID:     parts[0],
		Title:        parts[1],
		Duration:     dur,
		Artist:       parts[3],
		Description:  parts[4],
		ThumbnailURL: parts[5],
		OriginURL:    parts[6],
		AudioURL:     parts[7],
		SourceType:   "youtube",
		Exptime:      exptime,
	}
	return resp, nil
}

// ExtractExpireFromURL разбирает ссылку и возвращает значение параметра ?expire= в виде Unix-времени
func parseExptime(link string) (int64, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return 0, fmt.Errorf("некорректный URL: %w", err)
	}

	expStr := parsed.Query().Get("expire")
	if expStr == "" {
		return 0, fmt.Errorf("параметр 'expire' не найден в ссылке")
	}

	expTime, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать 'expire' в число: %w", err)
	}

	return expTime, nil
}

// IsExpiredSoon проверяет, истечёт ли ссылка раньше, чем закончится воспроизведение
func IsExpiredSoon(exptime int64, durationSeconds int) bool {
	now := time.Now().Unix()
	return now+int64(durationSeconds) > exptime
}
