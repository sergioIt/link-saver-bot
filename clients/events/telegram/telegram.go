package telegram

import "link-saver-bot/clients/telegram"

type Processor struct {

	tgClient *telegram.Client
	offset int
	// @TODO: add storage
}
