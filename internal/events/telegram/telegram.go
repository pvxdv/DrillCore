package telegram

import (
	tgClient "drillCore/internal/clients/telergam"
	"drillCore/internal/events"
	"drillCore/internal/model"
	debtStorage "drillCore/internal/storage/debt"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Processor struct {
	tg           *tgClient.Client
	offset       int
	storage      debtStorage.DebtStorage
	userSessions map[int]*SessionState
}

type Meta struct {
	ChatID int
	UserID int
}

// TODO extract to Redis

type SessionState struct {
	Action    string
	TempDebt  *model.Debt
	MessageID int
}

const (
	waitDescription = "awaiting_description"
	waitAmount      = "awaiting_amount"
	waitDate        = "awaiting_date"
)

const (
	callbackMainMenu = "main_menu"
	callbackStart    = "start"
	callbackListDebs = "list_debts"
	callbackAddDept  = "add_debt"
	callbackDelete   = "delete_"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
	ErrNoUpdatesFound   = errors.New("no updates found")
)

func New(client *tgClient.Client, storage debtStorage.DebtStorage) *Processor {
	//TODO extract to Redis
	sessions := make(map[int]*SessionState)
	return &Processor{
		tg:           client,
		storage:      storage,
		userSessions: sessions,
	}
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
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.Callback:
		return p.processCallback(event)
	default:
		return fmt.Errorf("failed to process event:%w", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process meta: %w", err)
	}

	session, exists := p.userSessions[meta.ChatID]
	if !exists {
		return p.showMainMenu(meta.ChatID)
	}

	switch session.Action {
	case waitDescription:
		return p.handleDebtDescription(meta.ChatID, event.Text)
	case waitAmount:
		return p.handleDebtAmount(meta.ChatID, event.Text)
	case waitDate:
		return p.handleDebtDate(meta.ChatID, event.Text)
	default:
		return p.showMainMenu(meta.ChatID)
	}
}

func (p *Processor) processCallback(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process meta: %w", err)
	}

	switch {
	case event.Text == callbackMainMenu || event.Text == callbackStart:
		return p.showMainMenu(meta.ChatID)

	case event.Text == callbackListDebs:
		return p.listDebts(meta.ChatID, meta.UserID)

	case event.Text == callbackAddDept:
		return p.startAddDebt(meta.ChatID)

	case strings.HasPrefix(event.Text, callbackDelete):
		debtID, _ := strconv.ParseInt(strings.TrimPrefix(event.Text, callbackDelete), 10, 64)
		return p.deleteDebt(meta.ChatID, debtID)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("failed to process meta: %w", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd tgClient.Update) events.Event {
	updType := fetchType(upd)

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
		m.UserID = upd.CallbackQuery.Message.From.ID
	default:
	}
	res.Meta = m
	return res
}

func fetchType(upd tgClient.Update) events.Type {
	switch {
	case upd.Message != nil:
		return events.Message
	case upd.CallbackQuery != nil:
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
