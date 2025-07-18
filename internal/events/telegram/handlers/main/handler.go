package mainmenu

import (
	tgClient "drillCore/internal/clients/telergam"
	"drillCore/internal/events"
	"drillCore/internal/events/telegram"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

const (
	prefix = "main_"

	cmdMenu = "menu"

	cbDebt    = "debt_menu"
	cbWeather = "weather_menu"
	cbTask    = "tasks_menu"
)

type Handler struct {
	tg     *tgClient.Client
	logger *zap.SugaredLogger
}

func New(tg *tgClient.Client, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		tg:     tg,
		logger: logger,
	}
}

func (h *Handler) ID() string {
	return prefix
}

func (h *Handler) CanHandle(event events.Event) bool {
	return strings.HasPrefix(event.Text, prefix)
}

func (h *Handler) Handle(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("get meta: %w", err)
	}

	action := strings.TrimPrefix(event.Text, prefix)

	switch action {
	case cmdMenu:
		return h.tg.SendMessageWithKeyboard(meta.ChatID, msgMenu, h.mainKeyboard())
	default:
		return nil
	}
}

func meta(event events.Event) (telegram.Meta, error) {
	res, ok := event.Meta.(telegram.Meta)
	if !ok {
		return telegram.Meta{}, fmt.Errorf("failed to process meta: %w", telegram.ErrUnknownMetaType)
	}
	return res, nil
}
