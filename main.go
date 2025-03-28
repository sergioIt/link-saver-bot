package main

import (
	"flag"
	"github.com/joho/godotenv"
	"link-saver-bot/clients/telegram"
	event_consumer "link-saver-bot/consumer/event-consumer"
	tgEvent "link-saver-bot/events/telegram"
	"link-saver-bot/storage/files"
	"log"
	"os"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "user_data"
	batchSize   = 100
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err) // This is just a warning
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN not found in environment")
	}

	eventProcessor := tgEvent.New(
		telegram.New(tgBotHost, token),
		files.New(storagePath))

	log.Print("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}
}

func mustToken() string {

	token := flag.String("tg-bot-token",
		"",
		"token to access telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("empty telegram token")
	}

	return *token
}
