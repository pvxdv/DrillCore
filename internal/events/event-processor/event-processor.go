package eventprocessor

import (
	"context"
	"errors"
	"fmt"

	"drillCore/internal/bot"
	"drillCore/internal/events"

	"go.uber.org/zap"
)

type HandlerManager interface {
	HandleEvent(ctx context.Context, event *events.Event) error
}

type Processor struct {
	tg         *bot.Client
	offset     int
	handlerMng HandlerManager
	logger     *zap.SugaredLogger
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrNoUpdatesFound   = errors.New("no updates found")
	ErrInvalidCommand   = errors.New("invalid command")
)

func New(tg *bot.Client, hm HandlerManager, logger *zap.SugaredLogger) *Processor {
	p := &Processor{
		tg:         tg,
		logger:     logger,
		handlerMng: hm,
	}

	return p
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]*events.Event, error) {
	updates, err := p.tg.Updates(ctx, p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updates: %w", err)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no updates found :%w", ErrNoUpdatesFound)
	}

	res := make([]*events.Event, 0, len(updates))
	for _, u := range updates {
		e, err := p.event(u)
		if err != nil {
			p.logger.Warnw("failed to fetch event", "event", u, "error", err)
		}

		res = append(res, e)
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, e *events.Event) error {
	p.logger.Debugf("processing event: %s", e.Text)

	if e.Type == events.Unknown {
		return fmt.Errorf("failed to process event: %w", ErrUnknownEventType)
	}

	return p.handlerMng.HandleEvent(ctx, e)
}

func (p *Processor) event(upd bot.Update) (*events.Event, error) {
	p.logger.Debugf("processing update:%+v", upd)

	updType := p.fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	var m events.Meta
	switch updType {
	case events.Message:
		m.ChatID = upd.Message.Chat.ID
		m.UserID = upd.Message.From.ID
	case events.Callback:
		m.ChatID = upd.CallbackQuery.Message.Chat.ID
		m.UserID = upd.CallbackQuery.From.ID
	case events.Unknown:
		return nil, ErrUnknownEventType
	}

	res.Meta = &m

	p.logger.Debugf("fetch event:%+v", res)

	return &res, nil
}

func (p *Processor) fetchType(upd bot.Update) events.Type {
	switch {
	case upd.Message != nil:
		p.logger.Debugf("fetch message type:%+v", upd.Message)
		return events.Message

	case upd.CallbackQuery != nil:
		p.logger.Debugf("fetch callback type:%+v", upd.CallbackQuery)
		return events.Callback
	default:
		return events.Unknown
	}
}

func fetchText(upd bot.Update) string {
	switch {
	case upd.Message != nil:

		return upd.Message.Text
	case upd.CallbackQuery != nil:
		return upd.CallbackQuery.Data
	default:
		return ""
	}
}
