package command

import (
	"context"
	"fmt"

	"drillCore/internal/bot"
	"drillCore/internal/events"
	"drillCore/internal/events/event-processor"
	"drillCore/internal/events/event-processor/manager"
	"drillCore/internal/session"

	"go.uber.org/zap"
)

type SessionManager interface {
	Get(ctx context.Context, userID int) (*session.Session, bool)
	Set(ctx context.Context, userID int, s *session.Session) error
	Delete(ctx context.Context, userID int) error
}

type Handler struct {
	tg     *bot.Client
	sesMng SessionManager
	logger *zap.SugaredLogger

	mainKB bot.ReplyMarkup
}

func New(tg *bot.Client, sm SessionManager, logger *zap.SugaredLogger) *Handler {
	h := &Handler{
		tg:     tg,
		sesMng: sm,
		logger: logger,
	}

	kb, err := h.mainKeyboard()
	if err != nil {
		h.logger.Fatal(err)
	}

	h.mainKB = kb

	return h
}

func (h *Handler) Type() manager.TypeHandler {
	return manager.CMDHandler
}

func (h *Handler) Handle(ctx context.Context, e *events.Event) error {
	defer func() { _ = h.sesMng.Delete(ctx, e.Meta.UserID) }()

	h.logger.Debugw("handling event in ", "handler", manager.CMDHandler, "event", e)

	if e.Type != events.Message {
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.InvalidEventType,
			),
		)
	}

	cmd, ok := manager.ParseCommand(e.Text)
	if !ok {
		return eventprocessor.ErrInvalidCommand
	}

	switch cmd {
	case manager.Start:
		return h.sendStartWithPhoto(ctx, e.Meta.ChatID, h.mainKB)

	case manager.Help:
		return h.tg.SendMessageWithKeyboard(ctx, e.Meta.ChatID, manager.MsgCMDHelp, h.mainKB)

	case manager.Debt:
		return h.tg.SendMessageWithKeyboard(ctx, e.Meta.ChatID, manager.MsgCMDDebt, h.mainKB)

	case manager.Recipe:
		return h.tg.SendMessageWithKeyboard(ctx, e.Meta.ChatID, manager.MsgCMDRecipe, h.mainKB)

	case manager.Gym:
		return h.tg.SendMessageWithKeyboard(ctx, e.Meta.ChatID, manager.MsgCMDGym, h.mainKB)

	case manager.Task:
		return h.tg.SendMessageWithKeyboard(ctx, e.Meta.ChatID, manager.MsgCMDTask, h.mainKB)

	default:
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.InvalidCommand,
			h.mainKB,
		)
	}
}

func (h *Handler) sendStartWithPhoto(ctx context.Context, chatID int, kb bot.ReplyMarkup) error {
	photoPath := "resources/static/welcome.jpg"

	err := h.tg.SendPhotoWithKeyBoard(ctx, chatID, photoPath, manager.MsgCMDStart, h.mainKB)
	if err != nil {
		h.logger.Errorf("failed to send photo: %v, falling back to text", err)

		return h.tg.SendMessageWithKeyboard(ctx, chatID, manager.MsgCMDStart, kb)
	}

	return nil
}

func (h *Handler) mainKeyboard() (bot.ReplyMarkup, error) {
	mainMenu, err := manager.CreateCallBack(manager.MainMenuHandler, manager.StepStart, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{{Text: manager.MainMenuButtonGeneral, CallbackData: mainMenu}},
	}), nil
}
