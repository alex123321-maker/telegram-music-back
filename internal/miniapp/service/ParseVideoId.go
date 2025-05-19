package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Регулярное выражение для ID YouTube
var idRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

// Главная функция извлечения ID
func ExtractYouTubeID(input string) (string, error) {
	// 1. Прямая проверка — это уже ID?
	if idRegex.MatchString(input) {
		return input, nil
	}

	// Попробуем распарсить URL
	u, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга URL: %w", err)
	}

	// Обработка различных хостов
	switch u.Host {
	case "youtu.be":
		// Короткий формат: https://youtu.be/ID
		return strings.TrimPrefix(u.Path, "/"), nil

	case "www.youtube.com", "youtube.com", "m.youtube.com":
		if strings.HasPrefix(u.Path, "/watch") {
			// Основной формат: https://www.youtube.com/watch?v=ID
			query := u.Query()
			return query.Get("v"), nil

		} else if strings.HasPrefix(u.Path, "/embed/") {
			// Embed формат: https://www.youtube.com/embed/ID
			return strings.TrimPrefix(u.Path, "/embed/"), nil

		} else if strings.HasPrefix(u.Path, "/shorts/") {
			// Shorts формат: https://www.youtube.com/shorts/ID
			return strings.TrimPrefix(u.Path, "/shorts/"), nil
		}
	}

	return "", fmt.Errorf("не удалось извлечь видео ID")
}
