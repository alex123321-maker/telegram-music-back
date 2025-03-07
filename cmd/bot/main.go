package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sourcecraft.dev/mainaccsteam035/telegram-music/internal/bot"
)

func main() {
	// Загрузка конфигурации
	cfg, err := bot.LoadConfig("configs/bot.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация бота
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatalf("Bot init failed: %v", err)
	}
	botAPI.Debug = cfg.Debug

	log.Printf("Authorized as @%s", botAPI.Self.UserName)

	// Настройка обработки обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case update := <-updates:
			handleUpdate(botAPI, update)
		case <-ctx.Done():
			log.Println("Shutting down...")
			botAPI.StopReceivingUpdates()
			return
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	
	switch update.Message.Command() {
	case "start":
		msg.Text = "🎵 Welcome to Music Bot!\n\nAvailable commands:\n/search - Find tracks\n/favorites - Your saved tracks"
	case "search":
		msg.Text = "🔍 Enter search query:"
	case "help":
		msg.Text = "ℹ️ Bot commands:\n/start - Initial setup\n/search - Find music\n/favorites - Saved tracks"
	default:
		msg.Text = "❌ Unknown command"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}