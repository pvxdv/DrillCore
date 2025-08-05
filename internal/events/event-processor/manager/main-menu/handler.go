package mainmenu

import (
	"context"
	"fmt"

	"drillCore/internal/bot"
	"drillCore/internal/events"
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
}

func New(tg *bot.Client, sm SessionManager, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		tg:     tg,
		sesMng: sm,
		logger: logger,
	}
}

func (h *Handler) Type() manager.TypeHandler {
	return manager.MainMenuHandler
}

func (h *Handler) Handle(ctx context.Context, e *events.Event) error {
	h.logger.Debugw("handling event in ", "handler", manager.MainMenuHandler, "event", e)

	if e.Type != events.Callback {
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.InvalidEventType,
			),
		)
	}

	cb, err := manager.ParseCallBack(e.Text)
	if err != nil {
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToGetCallBack,
		)
	}

	kb, err := h.mainKeyboard()
	if err != nil {
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	switch cb.Step {
	case manager.StepStart:
		return h.tg.SendMessageWithKeyboard(ctx, e.Meta.ChatID, manager.MsgMainMenu, kb)

	default:
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.InvalidStep,
				cb.Step,
			),
		)
	}
}

func (h *Handler) mainKeyboard() (bot.ReplyMarkup, error) {
	debtMenu, err := manager.CreateCallBack(manager.DebtHandler, manager.StepStart, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	ignore, err := manager.CreateCallBack(manager.IgnoreHandler, manager.StepIgnore, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{{Text: manager.DebtModuleButton, CallbackData: debtMenu}},
		{{Text: manager.RecipeModuleButton, CallbackData: ignore}},
		{{Text: manager.GymModuleButton, CallbackData: ignore}},
		{{Text: manager.TasksModuleButton, CallbackData: ignore}},
	}), nil
}
