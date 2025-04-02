package telegram

import (
	// "link-saver-bot/clients/events/"
	"errors"
	"link-saver-bot/lib/e"
	"link-saver-bot/storage"
	"log/slog"
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
	slog.Info("Received command",
		"command", commandText,
		"user", userName,
		"chatId", chatId)

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
		slog.Info("Unknown command received",
			"command", commandText,
			"user", userName)
		return p.tgClient.SendMessage(chatId, msgCommandUnknown)
	}
}

func isAddCmd(command string) bool {
	return isUrl(command)
}

func isUrl(command string) bool {
	u, err := url.Parse(command)
	return err == nil && u.Host != ""
}

func (p *Processor) savePage(chatID int, pageURL string, userName string) (err error) {
	slog.Info("Saving page",
		"url", pageURL,
		"user", userName,
		"chatId", chatID)

	page := &storage.Page{
		URL:      pageURL,
		UserName: userName,
	}

	isExists, err := p.storage.Exists(page)
	if err != nil {
		slog.Error("Failed to check page existence",
			"url", pageURL,
			"user", userName,
			"error", err)
		return e.Wrap("can't check if page exists", err)
	}

	if isExists {
		slog.Info("Page already exists",
			"url", pageURL,
			"user", userName)
		return p.tgClient.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		slog.Error("Failed to save page",
			"url", pageURL,
			"user", userName,
			"error", err)
		return err
	}

	slog.Info("Page saved successfully",
		"url", pageURL,
		"user", userName)
	return p.tgClient.SendMessage(chatID, msgSaved)
}

func (p *Processor) sendRandom(chatID int, userName string) (err error) {
	slog.Info("Fetching random page", "user", userName, "chatId", chatID)

	page, err := p.storage.PickRandom(userName)
	if err != nil {
		if errors.Is(err, storage.ErrNoSavedPages) {
			slog.Info("No saved pages found", "user", userName)
			return p.tgClient.SendMessage(chatID, msgNoSavedPages)
		}
		slog.Error("Failed to get random page",
			"user", userName,
			"error", err)
		return err
	}

	if err := p.tgClient.SendMessage(chatID, page.URL); err != nil {
		slog.Error("Failed to send random page",
			"url", page.URL,
			"user", userName,
			"error", err)
		return err
	}

	slog.Info("Random page sent and removed",
		"url", page.URL,
		"user", userName)
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	slog.Info("Sending help message", "chatId", chatID)
	return p.tgClient.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	slog.Info("Sending welcome message", "chatId", chatID)
	return p.tgClient.SendMessage(chatID, msgHello)
}
