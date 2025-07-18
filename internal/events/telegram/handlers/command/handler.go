package command

import (
	tgClient "drillCore/internal/clients/telergam"
	"drillCore/internal/events"
	"drillCore/internal/events/telegram"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

const (
	prefix = "/"

	cmdStart = "start"
	cmdHelp  = "help"

	cbMainMenu = "main_menu"
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
		return fmt.Errorf("failed to get meta: %w", err)
	}

	cmd := strings.TrimPrefix(event.Text, prefix)
	switch cmd {
	case cmdStart:
		return h.tg.SendMessageWithKeyboard(meta.ChatID, msgStart, h.mainKeyboard())
	case cmdHelp:
		return h.tg.SendMessageWithKeyboard(meta.ChatID, msgHelp, h.mainKeyboard())
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
