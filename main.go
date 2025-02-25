package main

import (
	"flag"
	"link-saver-bot/clients/telegram"
	event_consumer "link-saver-bot/consumer/event-consumer"
	tgEvent "link-saver-bot/events/telegram"
	"link-saver-bot/storage/files"
	"log"
)

const (
	tgBotHost   = "https://api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

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

	token := flag.String("token-bot-token", "", "token to access telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("empty telegram token")
	}

	return *token
}
