package telegram

import (
	// "link-saver-bot/clients/events/"
	"errors"
	"link-saver-bot/clients/telegram"
	"link-saver-bot/lib/e"
	"link-saver-bot/storage"
	"log"
	"net/url"
	"strings"
)

const (
	// add page is default command
	RndCmd   = "/rnd"   // get random stored link
	HelpCmd  = "/help"  // how to use this bot
	StartCmd = "/start" // hello from bot and help command
)

func (p *Processor) doCmd(commandText string, chatId int, userName string) error {

	commandText = strings.TrimSpace(commandText)
	log.Printf("got new command %s from user %s", commandText, userName)

	if isAddCmd(commandText) {
		return p.savePage(chatId, commandText, userName)
	}

	switch commandText {
	case RndCmd:
		return p.sendRandom(chatId, userName)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tgClient.SendMessage(chatId, msgCommandUnknown)
	}
}

func isAddCmd(command string) bool {
	return isUrl(command)

}

func isUrl(command string) bool {

	url, err := url.Parse(command)

	return err == nil && url.Host != ""
}

func (p *Processor) savePage(chatID int, pageURL string, userName string) (err error) {

	defer func() {
		err = e.Wrap("can't do save page command", err)
	}()

	sendMessage := NewMessageSender(chatID, p.tgClient)

	page := &storage.Page{
		URL:      pageURL,
		UserName: userName,
	}

	isExists, err := p.storage.Exists(page)

	if err != nil {
		return e.Wrap("can't check if page exists", err)
	}

	if isExists {
		err := sendMessage(msgAlreadyExists)
		if err != nil {
			return err
		}
		// return p.tgClient.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMessage(msgSaved); err != nil {
		return err
	}

	return nil
}

func NewMessageSender(chatId int, tg *telegram.Client) func(string) error {

	return func(message string) error {
		return tg.SendMessage(chatId, message)
	}
}

func (p *Processor) sendRandom(chatID int, userName string) (err error) {

	page, err := p.storage.PickRandom(userName)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tgClient.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tgClient.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tgClient.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tgClient.SendMessage(chatID, msgHello)
}
