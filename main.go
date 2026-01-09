package main

import (
	"flag"
	tgClient "links-bot/clients/telegram"
	event_consumer "links-bot/consumer/event-consumer"
	"links-bot/events/telegram"
	"links-bot/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.NewProcessor(
		tgClient.NewClient(tgBotHost, mustToken()),
		files.NewStorage(storagePath),
	)

	log.Println("servise started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("serviseis stoped: ", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
