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

	cbDebt     = "debt_menu"
	cbWeather  = "weather_menu"
	cbTask     = "tasks_menu"
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
	meta, err := telegram.ExtractMeta(event)
	if err != nil {
		return fmt.Errorf("failed to extract meta: %w", err)
	}

	action := strings.TrimPrefix(event.Text, prefix)

	switch action {
	case cmdMenu:
		return h.tg.SendMessageWithKeyboard(meta.ChatID, msgMenu, h.mainKeyboard())
	default:
		return nil
	}
}

func (h *Handler) mainKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: debtButton, CallbackData: cbDebt},
		},
		// [SPIRAL SYSTEMS CHARGING...]
		//{
		//    {Text: weatherButton, CallbackData: cbWeather},
		//    {Text: taskButton, CallbackData: cbTask},
		//},
		{
			{Text: MainConsoleButton, CallbackData: cbMainMenu},
		},
	})
}

const (
	debtButton        = "💢 DEPLOY DEBT DRILL SYSTEM (v2.1) 💢"
	weatherButton     = "🌪️ DEPLOY CLIMATE DRILL SYSTEM 🌪️"
	taskButton        = "🎯 DEPLOY TASK DRILLER 9000 🎯"
	MainConsoleButton = "🌀 SPIRAL COMMAND CENTER 🌀"
	msgMenu           = "🌀 SPIRAL COMMAND CENTER 🌀\n\n" +
		"AWAITING ORDERS, DRILL COMMANDER!\n\n" +
		"SELECT COMBAT PROTOCOL:\n"
)
