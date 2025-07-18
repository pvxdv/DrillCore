package telegram

import (
	tgClient "drillCore/internal/clients/telergam"
	"drillCore/internal/events"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

type Handler interface {
	CanHandle(event events.Event) bool
	Handle(event events.Event) error
	ID() string
}

type Processor struct {
	tg       *tgClient.Client
	offset   int
	logger   *zap.SugaredLogger
	handlers map[string]Handler
	sessions *events.SessionManager
}

type Meta struct {
	ChatID int
	UserID int
}

var (
	ErrNoHandlerFound   = errors.New("no handler found")
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
	ErrNoUpdatesFound   = errors.New("no updates found")
	ErrInvalidCommand   = errors.New("invalid command")
)

func New(tg *tgClient.Client, sm *events.SessionManager, logger *zap.SugaredLogger, handlers ...Handler) *Processor {
	p := &Processor{
		tg:       tg,
		logger:   logger,
		handlers: make(map[string]Handler),
		sessions: sm,
	}

	for _, h := range handlers {
		p.handlers[h.ID()] = h
	}

	return p
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updates: %w", err)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no updates found :%w", ErrNoUpdatesFound)
	}

	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, p.event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	if event.Type == events.Unknown {
		return fmt.Errorf("failed to process event :%w", ErrUnknownEventType)
	}

	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	if session, exists := p.sessions.Get(meta.UserID); exists {
		for _, h := range p.handlers {
			if h.ID() == session.HandlerID {
				p.logger.Debugw("handle event with handler", "event", event.Text, "handlerId", h.ID())
				return h.Handle(event)
			}
		}

		p.sessions.Delete(meta.UserID)
	}

	for _, h := range p.handlers {
		if h.CanHandle(event) {
			p.logger.Debugw("handle event with handler", "event", event.Text, "handlerId", h.ID())
			return h.Handle(event)
		}
	}

	return fmt.Errorf("failed to process event:%s :%w", event.Text, ErrNoHandlerFound)
}

func (p *Processor) event(upd tgClient.Update) events.Event {
	p.logger.Debugf("resived update:%+v", upd)
	updType := p.fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	var m Meta
	switch updType {
	case events.Message:
		m.ChatID = upd.Message.Chat.ID
		m.UserID = upd.Message.From.ID
	case events.Callback:
		m.ChatID = upd.CallbackQuery.Message.Chat.ID
		m.UserID = upd.CallbackQuery.From.ID
	case events.Unknown:
		p.logger.Warnf("unknown event:%+v", upd)
	}
	res.Meta = m
	p.logger.Debugf("return event: %+v", res)
	return res
}

func (p *Processor) fetchType(upd tgClient.Update) events.Type {
	switch {
	case upd.Message != nil:
		p.logger.Debugf("got message: %+v", upd.Message)
		return events.Message
	case upd.CallbackQuery != nil:
		p.logger.Debugf("got callback query: %+v", upd.CallbackQuery)
		return events.Callback
	default:
		return events.Unknown
	}
}

func fetchText(upd tgClient.Update) string {
	switch {
	case upd.Message != nil:
		return upd.Message.Text
	}
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}

	return ""
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("failed to process meta: %w", ErrUnknownMetaType)
	}
	return res, nil
}
