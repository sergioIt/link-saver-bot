package telegram

import (
	"errors"
	events "link-saver-bot/clients/events"
	"link-saver-bot/clients/telegram"
	"link-saver-bot/lib/e"
	"link-saver-bot/storage"
)

type Processor struct {
	tgClient *telegram.Client
	offset   int
	storage  storage.Storage
}

type Meta struct {
	ChatId   int
	UserName string
}

var (
	errUnknownEventType = errors.New("unknown event type")
	errMetaNotFound     = errors.New("meta not found")
)

func New(tgClient *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tgClient: tgClient,
		storage:  storage,
	}
}

func (p *Processor) Fetch(limit, offset int) ([]events.Event, error) {
	updates, err := p.tgClient.Updates(limit, p.offset)
	if err != nil {
		return nil, e.Wrap("can't get events, error: ", err)
	}
	// return empty result if no updates were fetched
	if len(updates) == 0 {
		return nil, nil
	}

	result := make([]events.Event, 0, len(updates))

	// transform updates into events
	for _, u := range updates {
		result = append(result, makeEvent(u))
	}

	p.offset = updates[len(updates)-1].ID + 1 // this is for the next fetch
	return result, nil
}

// process event depending of its type:
// process event depending of its type:
// 1. If message -> process message
// 2. If unknown -> return error
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", errUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := getMeta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatId, meta.UserName); err != nil {
		return e.Wrap("can't process command", err)
	}

	return nil
}

func getMeta(event events.Event) (Meta, error) {

	meta, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", errMetaNotFound)
	}

	return meta, nil
}

func makeEvent(update telegram.Updates) events.Event {

	updateType := fetchType(update)

	res := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatId:   update.Message.Chat.ID,
			UserName: update.Message.From.UserName,
		}
	}

	return res
}

func fetchType(update telegram.Updates) events.Type {

	if update.Message == nil {
		return events.Unknown
	}

	return events.Message

}

func fetchText(update telegram.Updates) string {

	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}
