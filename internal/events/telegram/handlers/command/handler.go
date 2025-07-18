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

	cmdStart = "/start"
	cmdHelp  = "/help"

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

	switch event.Text {
	case cmdStart:
		return h.tg.SendMessageWithKeyboard(meta.ChatID, msgStart, h.mainKeyboard())
	case cmdHelp:
		return h.tg.SendMessageWithKeyboard(meta.ChatID, msgHelp, h.mainKeyboard())
	default:
		return nil
	}
}

const (
	msgStart = `🌀 PIERCE THE HEAVENS WITH YOUR SPIRAL POWER! 🌀

YOUR DRILL IS THE DRILL THAT WILL CREATE TOMORROW!
YOUR DRILL IS THE DRILL THAT SMASHES LIMITS!
YOUR DRILL IS... this productivity bot!

⚙️ CURRENT DRILL SYSTEMS:
• DEBT DRILL SYSTEM (v2.1) [OPERATIONAL]
• LAZINESS ANNIHILATOR [STANDBY]
• PROCRASTINATION CRUSHER [STANDBY]

⚡ COMING WHEN WE PIERCE THE DEVELOPMENT HELL:
• TASK BUSTER 9000
• WEATHER IMPACT SYSTEM
• ANTI-PROCASTINATION DRILL

DEPLOY MENU BELOW TO BEGIN DRILLING!`

	msgHelp = `💢 SPIRAL COMMAND TRANSMISSION 💢

🌀 /start - ACTIVATE MAIN DRILL
📡 /help - DISPLAY COMBAT MANUAL

THE POWER TO CHANGE YOUR LIFE IS THE POWER TO DRILL!
WHO THE HELL DO YOU THINK YOU ARE?!`

	mainButton = "🌀 DEPLOY SPIRAL COMMAND CENTER 🌀"
)

func (h *Handler) mainKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{{Text: mainButton, CallbackData: cbMainMenu}},
	})
}
