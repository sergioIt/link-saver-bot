package main

import (
	"flag"
	// "github.com/joho/godotenv"
	"link-saver-bot/clients/telegram"
	event_consumer "link-saver-bot/consumer/event-consumer"
	tgEvent "link-saver-bot/events/telegram"
	"link-saver-bot/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	//token: 7655555780:AAE3tTP1MwCXbbSdz7D_4FzPARkgb-LO4yY

	// tgClient := telegram.New(tgBotHost, mustToken())

	eventProcessor := tgEvent.New(telegram.New(tgBotHost, mustToken()),
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
