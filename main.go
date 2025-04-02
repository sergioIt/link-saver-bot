package main

import (
	"github.com/joho/godotenv"
	"link-saver-bot/clients/telegram"
	event_consumer "link-saver-bot/consumer/event-consumer"
	tgEvent "link-saver-bot/events/telegram"
	"link-saver-bot/storage/files"
	"log/slog"
	"os"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "user_data"
	batchSize   = 100
)

func main() {
	// Setup structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("Error loading .env file", "error", err)
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		slog.Error("TELEGRAM_TOKEN not found in environment")
		os.Exit(1)
	}

	eventProcessor := tgEvent.New(
		telegram.New(tgBotHost, token),
		files.New(storagePath))

	slog.Info("Service started", "host", tgBotHost, "storage", storagePath)

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		slog.Error("Service error", "error", err)
		os.Exit(1)
	}
}
