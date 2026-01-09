package telegram

import (
	"errors"
	"links-bot/clients/telegram"
	"links-bot/events"
	e "links-bot/lib"
	"links-bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func NewProcessor(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't  get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, upd := range updates {
		res = append(res, event(upd))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return ErrUnknownEventType
	}

}

func (p Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMetaType
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.UserName,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
