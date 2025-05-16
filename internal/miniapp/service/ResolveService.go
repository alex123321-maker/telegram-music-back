package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	AudioURL     string `json:"audio_url"`
}

func ResolveUrl(req string) (ResolveMediaResponse, error) {
	dl := ytdlp.New().
		Quiet().
		NoWarnings().
		NoPlaylist().
		Format("bestaudio[ext=m4a]/bestaudio").
		Print("%(id)s␞%(title)s␞%(duration)s␞%(uploader)s␞%(description)s␞%(thumbnail)s␞%(webpage_url)s␞%(url)s").
		Proxy("socks5://127.0.0.1:10808")

	raw, err := dl.Run(context.TODO(), req)
	if err != nil {
		return ResolveMediaResponse{}, fmt.Errorf("ошибка при полчении ответа от dlp")
	}

	parts := strings.Split(string(raw.Stdout), "␞")
	if len(parts) != 8 {
		return ResolveMediaResponse{}, fmt.Errorf("yt-dlp вернул неверный формат, %d", len(parts))
	}

	dur, _ := strconv.Atoi(parts[2])

	resp := ResolveMediaResponse{
		SourceID:     parts[0],
		Title:        parts[1],
		Duration:     dur,
		Artist:       parts[3],
		Description:  parts[4],
		ThumbnailURL: parts[5],
		AudioURL:     parts[7],
		SourceType:   "youtube",
	}
	return resp, nil
}
