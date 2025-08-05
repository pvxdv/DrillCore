package debt

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"drillCore/internal/bot"
	"drillCore/internal/events"
	"drillCore/internal/events/event-processor/manager"
	"drillCore/internal/model"
	"drillCore/internal/session"

	"go.uber.org/zap"
)

type SessionManager interface {
	Get(ctx context.Context, userID int) (*session.Session, bool)
	Set(ctx context.Context, userID int, s *session.Session) error
	Delete(ctx context.Context, userID int) error
}

type Storage interface {
	Save(ctx context.Context, debt *model.Debt) (int64, error)
	Debts(ctx context.Context, userID int64) ([]*model.Debt, error)
	Update(ctx context.Context, debt *model.Debt) error
	Delete(ctx context.Context, id int64) error
	Debt(ctx context.Context, id int64) (*model.Debt, error)
}

type Handler struct {
	tg      *bot.Client
	sesMng  SessionManager
	storage Storage
	logger  *zap.SugaredLogger

	menuKeyBoard     bot.ReplyMarkup
	cancelKeyBoard   bot.ReplyMarkup
	editMenuKeyBoard bot.ReplyMarkup
}

func New(tg *bot.Client, sm SessionManager, storage Storage, logger *zap.SugaredLogger) *Handler {
	h := &Handler{
		tg:      tg,
		sesMng:  sm,
		storage: storage,
		logger:  logger,
	}

	menuKeyboard, err := h.menuKeyboard()
	if err != nil {
		h.logger.Fatal(err)
	}

	cancelKeyboard, err := h.cancelKeyboard()
	if err != nil {
		h.logger.Fatal(err)
	}

	editMenuKeyBoard, err := h.editKeyboard()
	if err != nil {
		h.logger.Fatal(err)
	}

	h.cancelKeyBoard = cancelKeyboard
	h.menuKeyBoard = menuKeyboard
	h.editMenuKeyBoard = editMenuKeyBoard

	return h
}

func (h *Handler) Type() manager.TypeHandler {
	return manager.DebtHandler
}

func (h *Handler) Handle(ctx context.Context, e *events.Event) error {
	h.logger.Debugw(
		"handling event in ",
		"handler", manager.DebtHandler,
		"event", e,
	)

	switch e.Type {
	case events.Message:
		return h.handleMessage(ctx, e)

	case events.Callback:
		cb, err := manager.ParseCallBack(e.Text)
		if err != nil {
			h.logger.Error(err)

			h.cleanupSession(ctx, e.Meta.UserID)

			return err
		}

		return h.handleCallBack(ctx, cb, e.Meta)

	default:
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.InvalidEventType,
			),
		)
	}
}

func (h *Handler) handleMessage(ctx context.Context, e *events.Event) error {
	ses, state, err := h.getSession(ctx, e.Meta.UserID, e.Meta.ChatID)
	if err != nil {
		return fmt.Errorf("failed to handle message for userID:%d :%v", e.Meta.UserID, err)
	}

	switch state.Step {
	case manager.StepAddDescription:
		return h.addDescription(ctx, e, ses, state)

	case manager.StepAddAmount:
		return h.addAmount(ctx, e, ses, state)

	case manager.StepPayAmount:
		return h.payAmount(ctx, e, ses, state)

	case manager.StepEditDescription:
		return h.editDescription(ctx, e, ses, state)

	case manager.StepEditAmount:
		return h.editAmount(ctx, e, ses, state)

	default:
		h.logger.Errorf("failed to handle event: %v for user %d", e, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.InvalidStep,
				state.Step,
			),
		)
	}
}

func (h *Handler) handleCallBack(ctx context.Context, cb *manager.CallBack, meta *events.Meta) error {
	h.logger.Debugf("handling callback event: %v", cb)

	switch cb.Step {
	case manager.StepStart:
		return h.debtStart(ctx, meta.ChatID, meta.UserID)

	case manager.StepList:
		return h.list(ctx, meta.ChatID, meta.UserID)

	case manager.StepAddStart:
		return h.addStart(ctx, meta.ChatID, meta.UserID)

	case manager.StepAddDescription:
		return h.addDescriptionRedirect(ctx, meta.ChatID, meta.UserID)

	case manager.StepAddAmount:
		return h.addAmountRedirect(ctx, meta.ChatID, meta.UserID)

	case manager.StepAddFinish:
		return h.addFinish(ctx, meta.ChatID, meta.UserID)

	case manager.StepDeleteStart:
		return h.beforeSelect(ctx, meta.ChatID, meta.UserID, manager.DebtHandler, manager.StepDeleteStart, manager.DebtHandler, manager.StepDeleteConfirm, manager.MsgDeleteStart)

	case manager.StepDeleteConfirm:
		return h.deleteConfirm(ctx, meta.ChatID, meta.UserID)

	case manager.StepDeleteFinish:
		return h.deleteFinish(ctx, meta.ChatID, meta.UserID)

	case manager.StepSelect:
		return h.selectDebt(ctx, meta.ChatID, meta.UserID, cb.Data)

	case manager.StepEditStart:
		return h.beforeSelect(ctx, meta.ChatID, meta.UserID, manager.DebtHandler, manager.StepEditStart, manager.DebtHandler, manager.StepEditMenu, manager.MsgEditStart)

	case manager.StepEditMenu:
		return h.editMenu(ctx, meta.ChatID, meta.UserID)

	case manager.StepEnterAmount:
		return h.enterAmount(ctx, meta.ChatID, meta.UserID)

	case manager.StepEnterDescription:
		return h.enterDescription(ctx, meta.ChatID, meta.UserID)

	case manager.StepEnterDate:
		return h.enterDate(ctx, meta.ChatID, meta.UserID)

	case manager.StepEditDate:
		return h.editDate(ctx, meta.ChatID, meta.UserID)

	case manager.StepEditFinish:
		return h.editFinish(ctx, meta.ChatID, meta.UserID)

	case manager.StepPayStart:
		return h.beforeSelect(ctx, meta.ChatID, meta.UserID, manager.DebtHandler, manager.StepPayStart, manager.DebtHandler, manager.StepEnterPayment, manager.MsgPayStart)

	case manager.StepEnterPayment:
		return h.enterPayment(ctx, meta.ChatID, meta.UserID)

	case manager.StepPayFinish:
		return h.payFinish(ctx, meta.ChatID, meta.UserID)

	default:
		h.logger.Errorf("failed to handle call back: %v for user %d", cb, meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			meta.ChatID,
			fmt.Sprintf(
				manager.InvalidStep,
				cb.Step,
			),
		)
	}
}

func (h *Handler) addStart(ctx context.Context, chatID, userID int) error {
	d := &model.Debt{
		UserID: int64(userID),
	}

	st := &manager.State{
		Handler:  manager.DebtHandler,
		Step:     manager.StepAddDescription,
		TempDebt: d,
	}

	s := &session.Session{
		State: st,
	}

	err := h.sesMng.Set(ctx, userID, s)
	if err != nil {
		h.logger.Errorf("failed to set session for user: %d", userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessage(
		ctx,
		chatID,
		manager.MsgAddDescription,
	)
}

func (h *Handler) addDescription(ctx context.Context, e *events.Event, ses *session.Session, state *manager.State) error {
	if strings.TrimSpace(e.Text) == "" {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidDescriptionEmpty,
			h.cancelKeyBoard,
		)
	}

	if len(e.Text) > 1000 {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.MsgInvalidDescriptionLength,
				len(e.Text),
			),
			h.cancelKeyBoard,
		)
	}

	state.TempDebt.Description = e.Text
	state.Step = manager.StepAddAmount

	ses.State = state
	err := h.sesMng.Set(ctx, e.Meta.UserID, ses)
	if err != nil {
		h.logger.Errorf("failed to save session for user %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.UserID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessage(
		ctx,
		e.Meta.ChatID,
		fmt.Sprintf(
			manager.MsgAddAmount,
			strings.ToUpper(e.Text),
		),
	)
}

func (h *Handler) addAmount(ctx context.Context, e *events.Event, ses *session.Session, state *manager.State) error {
	if strings.TrimSpace(e.Text) == "" {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidAmountEmpty,
			h.cancelKeyBoard,
		)
	}

	amount, err := strconv.ParseInt(e.Text, 10, 64)
	if err != nil || amount <= 0 {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidAmountConvertErr,
			h.cancelKeyBoard,
		)
	}

	state.TempDebt.Amount = amount

	state.BackHandler = manager.DebtHandler
	state.BackStep = manager.StepAddAmount

	state.Handler = manager.DateHandler
	state.Step = manager.StepYear

	state.NextHandler = manager.DebtHandler
	state.NextStep = manager.StepAddFinish

	ses.State = state
	err = h.sesMng.Set(ctx, e.Meta.UserID, ses)
	if err != nil {
		h.logger.Errorf("failed to save session for user %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.UserID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	kb, err := h.dateKeyboard(manager.DebtHandler, manager.StepAddAmount)
	if err != nil {
		h.cleanupSession(ctx, e.Meta.UserID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	msg := fmt.Sprintf(manager.MsgAddDate, formatMoney(amount)) + manager.MsgStartDateFlow

	return h.tg.SendMessageWithKeyboard(
		ctx,
		e.Meta.ChatID,
		msg,
		kb,
	)
}

func (h *Handler) addDescriptionRedirect(ctx context.Context, chatID, userID int) error {
	s, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to finish pay for userID:%d :%v", userID, err)
	}

	state.Handler = manager.DebtHandler
	state.Step = manager.StepAddDescription

	s.State = state

	err = h.sesMng.Set(ctx, userID, s)
	if err != nil {
		h.logger.Errorf("failed to set session for user: %d", userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessage(
		ctx,
		chatID,
		manager.MsgAddDescription,
	)
}

func (h *Handler) addAmountRedirect(ctx context.Context, chatID, userID int) error {
	s, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to finish pay for userID:%d :%v", userID, err)
	}

	state.Handler = manager.DebtHandler
	state.Step = manager.StepAddAmount

	s.State = state

	err = h.sesMng.Set(ctx, userID, s)
	if err != nil {
		h.logger.Errorf("failed to set session for user: %d", userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessage(
		ctx,
		chatID,
		fmt.Sprintf(
			manager.MsgAddAmount,
			strings.ToUpper(state.TempDebt.Description),
		),
	)
}

func (h *Handler) addFinish(ctx context.Context, chatID, userID int) error {
	defer h.cleanupSession(ctx, userID)

	_, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to handle finish step for userID:%d :%v", userID, err)
	}

	if state.TempDate == nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgDateNotSet,
			h.menuKeyBoard,
		)
	}

	state.TempDebt.ReturnDate = state.TempDate

	_, err = h.storage.Save(ctx, state.TempDebt)
	if err != nil {
		h.logger.Errorf("failed to save debt for user:%d : %v", userID, err)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSaveDebt,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		fmt.Sprintf(
			manager.MsgSavedDebt,
			strings.ToUpper(state.TempDebt.Description),
			formatMoney(state.TempDebt.Amount),
			state.TempDebt.ReturnDate.Format("02.01.2006"),
		),
		h.menuKeyBoard,
	)
}

func (h *Handler) deleteConfirm(ctx context.Context, chatID, userID int) error {
	_, _, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to confirm delete for userID:%d :%v", userID, err)
	}

	kb, err := h.confirmKeyboard(manager.StepDeleteFinish)
	if err != nil {
		return h.tg.SendMessage(
			ctx,
			chatID,
			manager.FailedToCreateKeyboard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		manager.MsgConfirmDeleteWarning,
		kb,
	)
}

func (h *Handler) deleteFinish(ctx context.Context, chatID, userID int) error {
	defer h.cleanupSession(ctx, userID)

	_, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete finish for userID:%d :%v", userID, err)
	}

	if err := h.storage.Delete(ctx, state.TempDebt.ID); err != nil {
		h.logger.Errorf("failed to delete debt %d: %v", state.TempDebt.ID, err)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToDeleteDebt,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		fmt.Sprintf(
			manager.MsgDeleteDebt,
			strings.ToUpper(state.TempDebt.Description),
			formatMoney(state.TempDebt.Amount),
		),
		h.menuKeyBoard,
	)
}

func (h *Handler) enterPayment(ctx context.Context, chatID, userID int) error {
	ses, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to enter payment for userID:%d :%v", userID, err)
	}

	state.Handler = manager.DebtHandler
	state.Step = manager.StepPayAmount

	ses.State = state

	err = h.sesMng.Set(ctx, userID, ses)
	if err != nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessage(
		ctx,
		chatID,
		manager.MsgEnterPayment,
	)
}

func (h *Handler) payAmount(ctx context.Context, e *events.Event, ses *session.Session, state *manager.State) error {
	if strings.TrimSpace(e.Text) == "" {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidAmountEmpty,
			h.cancelKeyBoard,
		)
	}

	amount, err := strconv.ParseInt(e.Text, 10, 64)
	if err != nil || amount <= 0 {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidAmountConvertErr,
			h.cancelKeyBoard,
		)
	}

	if amount > state.TempDebt.Amount {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.MsgToLargeAmount,
				formatMoney(state.TempDebt.Amount),
			),
			h.cancelKeyBoard,
		)
	}

	newAmount := state.TempDebt.Amount - amount

	state.TempDebt.Amount = newAmount
	ses.State = state

	err = h.sesMng.Set(ctx, e.Meta.UserID, ses)
	if err != nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	confirmMsg := fmt.Sprintf(
		manager.MsgPayConfirm,
		strings.ToUpper(state.TempDebt.Description),
		formatMoney(state.TempDebt.Amount),
		formatMoney(amount),
		formatMoney(newAmount),
	)

	confirmKb, err := h.confirmKeyboard(manager.StepPayFinish)
	if err != nil {
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		e.Meta.ChatID,
		confirmMsg,
		confirmKb,
	)
}

func (h *Handler) payFinish(ctx context.Context, chatID, userID int) error {
	h.cleanupSession(ctx, userID)

	_, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to finish pay for userID:%d :%v", userID, err)
	}

	switch state.TempDebt.Amount {
	case 0:
		if err := h.storage.Delete(ctx, state.TempDebt.ID); err != nil {
			h.logger.Errorf("failed to delete debt %d: %v", state.TempDebt.ID, err)

			return h.tg.SendMessageWithKeyboard(
				ctx,
				chatID,
				manager.MsgFailedToDeleteDebt,
				h.menuKeyBoard,
			)
		}

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			fmt.Sprintf(
				manager.MsgPayToDelete,
				strings.ToUpper(state.TempDebt.Description),
			),
			h.menuKeyBoard,
		)

	default:
		if err := h.storage.Update(ctx, state.TempDebt); err != nil {
			h.logger.Errorf("failed to update debt %d: %v", state.TempDebt.ID, err)

			return h.tg.SendMessageWithKeyboard(
				ctx,
				chatID,
				manager.MsgFailedToUpdateDebt,
				h.menuKeyBoard,
			)
		}

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			fmt.Sprintf(
				manager.MsgPayToUpdate,
				strings.ToUpper(state.TempDebt.Description),
				formatMoney(state.TempDebt.Amount),
			),
			h.menuKeyBoard,
		)
	}
}

func (h *Handler) editMenu(ctx context.Context, chatID, userID int) error {
	_, _, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to show edit menu for userID:%d :%v", userID, err)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		manager.MsgEditMenu,
		h.editMenuKeyBoard,
	)
}

func (h *Handler) enterAmount(ctx context.Context, chatID, userID int) error {
	ses, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to enter amount for userID:%d :%v", userID, err)
	}

	state.Handler = manager.DebtHandler
	state.Step = manager.StepEditAmount

	ses.State = state

	err = h.sesMng.Set(ctx, userID, ses)
	if err != nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessage(
		ctx,
		chatID,
		manager.MsgEnterAmount,
	)
}

func (h *Handler) editAmount(ctx context.Context, e *events.Event, ses *session.Session, state *manager.State) error {
	if strings.TrimSpace(e.Text) == "" {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidAmountEmpty,
			h.cancelKeyBoard,
		)
	}

	amount, err := strconv.ParseInt(e.Text, 10, 64)
	if err != nil || amount <= 0 {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidAmountConvertErr,
			h.cancelKeyBoard,
		)
	}

	state.TempDebt.Amount = amount
	state.Step = manager.StepEditMenu
	ses.State = state

	err = h.sesMng.Set(ctx, e.Meta.UserID, ses)
	if err != nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		e.Meta.ChatID,
		fmt.Sprintf(
			manager.MsgEditAmount,
			formatMoney(amount),
		),
		h.editMenuKeyBoard,
	)
}

func (h *Handler) enterDescription(ctx context.Context, chatID, userID int) error {
	ses, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to enter description for userID:%d :%v", userID, err)
	}

	state.Handler = manager.DebtHandler
	state.Step = manager.StepEditDescription

	ses.State = state

	err = h.sesMng.Set(ctx, userID, ses)
	if err != nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessage(
		ctx,
		chatID,
		manager.MsgEnterDescription,
	)
}

func (h *Handler) editDescription(ctx context.Context, e *events.Event, ses *session.Session, state *manager.State) error {
	if strings.TrimSpace(e.Text) == "" {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidDescriptionEmpty,
			h.cancelKeyBoard,
		)
	}

	if len(e.Text) > 1000 {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.MsgInvalidDescriptionLength,
				len(e.Text),
			),
			h.cancelKeyBoard,
		)
	}

	state.TempDebt.Description = e.Text
	state.Step = manager.StepEditMenu
	ses.State = state

	err := h.sesMng.Set(ctx, e.Meta.UserID, ses)
	if err != nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		e.Meta.ChatID,
		fmt.Sprintf(
			manager.MsgEditDescription,
			strings.ToUpper(e.Text),
		),
		h.editMenuKeyBoard,
	)
}

func (h *Handler) enterDate(ctx context.Context, chatID, userID int) error {
	ses, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to edit date for userID:%d :%v", userID, err)
	}

	state.BackHandler = manager.DebtHandler
	state.BackStep = manager.StepEditMenu
	state.NextHandler = manager.DebtHandler
	state.NextStep = manager.StepEditDate
	state.Handler = manager.DateHandler
	state.Step = manager.StepYear

	ses.State = state
	err = h.sesMng.Set(ctx, userID, ses)
	if err != nil {
		h.logger.Errorf("failed to save session for user %d", chatID)

		h.cleanupSession(ctx, userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	kb, err := h.dateKeyboard(manager.DebtHandler, manager.StepEditMenu)
	if err != nil {
		h.cleanupSession(ctx, userID)

		return h.tg.SendMessage(
			ctx,
			chatID,
			manager.FailedToCreateKeyboard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		manager.MsgStartDateFlow,
		kb,
	)
}

func (h *Handler) editDate(ctx context.Context, chatID, userID int) error {
	ses, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to finish edit for userID:%d :%v", userID, err)
	}

	if state.TempDate == nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			fmt.Sprintf(
				manager.MsgDateNotSet,
			),
			h.editMenuKeyBoard,
		)
	}

	state.TempDebt.ReturnDate = state.TempDate
	ses.State = state

	err = h.sesMng.Set(ctx, userID, ses)
	if err != nil {
		h.logger.Errorf("failed to save session for user %d", chatID)

		h.cleanupSession(ctx, userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		fmt.Sprintf(
			manager.MsgEditDate,
			state.TempDate.Format("02.01.2006"),
		),
		h.editMenuKeyBoard,
	)
}

func (h *Handler) editFinish(ctx context.Context, chatID, userID int) error {
	defer h.cleanupSession(ctx, userID)

	_, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to finish edit for userID:%d :%v", userID, err)
	}

	err = h.storage.Update(ctx, state.TempDebt)
	if err != nil {
		h.logger.Errorf("failed to save debt for user:%d : %v", userID, err)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToUpdateDebt,
			h.menuKeyBoard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		fmt.Sprintf(
			manager.MsgFinishEdit,
			strings.ToUpper(state.TempDebt.Description),
			formatMoney(state.TempDebt.Amount),
		),
		h.menuKeyBoard,
	)
}

func (h *Handler) beforeSelect(ctx context.Context, chatID, userID int, backH manager.TypeHandler, backS manager.Step,
	nextH manager.TypeHandler, nextS manager.Step, msg string) error {
	debts, err := h.storage.Debts(ctx, int64(userID))
	if err != nil {
		h.logger.Errorf("failed to get debts: %v", err)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToGetDebt,
			h.menuKeyBoard,
		)
	}

	if len(debts) == 0 {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.NoDebtsPhrases[rand.Intn(len(manager.NoDebtsPhrases))],
			h.menuKeyBoard,
		)
	}

	st := &manager.State{
		BackHandler: backH,
		BackStep:    backS,
		Handler:     manager.DebtHandler,
		Step:        manager.StepSelect,
		NextHandler: nextH,
		NextStep:    nextS,
	}

	s := &session.Session{
		State: st,
	}

	err = h.sesMng.Set(ctx, userID, s)
	if err != nil {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	sDebts := sortDebts(debts)

	kb, err := h.selectKeyboard(sDebts)
	if err != nil {
		return h.tg.SendMessage(
			ctx,
			chatID,
			manager.FailedToCreateKeyboard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		msg,
		kb,
	)
}

func (h *Handler) selectDebt(ctx context.Context, chatID, userID int, data string) error {
	ses, state, err := h.getSession(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to select debt for userID:%d :%v", userID, err)
	}

	debtID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		h.logger.Errorf("failed to extract debt id from debt %s", data)

		h.cleanupSession(ctx, userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToExtractDebtId,
			h.menuKeyBoard,
		)
	}

	debt, err := h.storage.Debt(ctx, debtID)
	if err != nil {
		h.logger.Errorf("failed to get debt for user: %d :%v", userID, err)

		h.cleanupSession(ctx, userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToGetDebt,
			h.menuKeyBoard,
		)
	}

	if debt.UserID != int64(userID) {
		h.logger.Errorf("debtID:%d relate to user:%d, request user:%d", debt.ID, debt.UserID, userID)

		h.cleanupSession(ctx, userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgUserIdNotEqualDebtId,
			h.menuKeyBoard,
		)
	}

	state.TempDebt = debt
	ses.State = state

	err = h.sesMng.Set(ctx, userID, ses)
	if err != nil {
		h.logger.Errorf("failed to set state for user %d", userID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.MsgFailedToSetSession,
			h.menuKeyBoard,
		)
	}

	redirectKb, err := h.redirectKeyboard(state.BackHandler, state.BackStep, state.NextHandler, state.NextStep)
	if err != nil {
		return h.tg.SendMessage(
			ctx,
			chatID,
			manager.FailedToCreateKeyboard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		fmt.Sprintf(
			manager.MsgDebtSelected,
			debt.Description,
			formatMoney(debt.Amount),
			debtStatus(debt),
		),
		redirectKb,
	)
}

func (h *Handler) debtStart(ctx context.Context, chatID int, userID int) error {
	h.cleanupSession(ctx, userID)

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		manager.MsgMenu,
		h.menuKeyBoard,
	)
}

func (h *Handler) list(ctx context.Context, chatID, userID int) error {
	debts, err := h.storage.Debts(ctx, int64(userID))
	if err != nil {
		h.logger.Errorf("failed to get debts for user: %d : %v", userID, err)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			manager.FailedToGetDebts,
			h.menuKeyBoard,
		)
	}

	if len(debts) == 0 {
		var sb strings.Builder
		sb.WriteString(manager.SpiralDelimiter)
		sb.WriteString(manager.NoDebtsPhrases[rand.Intn(len(manager.NoDebtsPhrases))] + "\n\n")
		sb.WriteString(manager.SpiralDelimiter)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			chatID,
			sb.String(),
			h.menuKeyBoard,
		)
	}

	sortedDebts := sortDebts(debts)

	total := int64(0)
	for _, d := range sortedDebts {
		total += d.Amount
	}

	var sb strings.Builder
	sb.WriteString(manager.SpiralDelimiter)
	sb.WriteString(manager.DebtTitles[rand.Intn(len(manager.DebtTitles))] + "\n\n")
	sb.WriteString(manager.SpiralDelimiter)

	for i, debt := range sortedDebts {
		marker := manager.RageEmoji
		if i%2 == 0 {
			marker = manager.SpiralEmoji
		}
		if debt.ReturnDate != nil && debt.ReturnDate.Before(time.Now()) {
			marker = manager.SkullEmoji
		}

		sb.WriteString(
			fmt.Sprintf(
				manager.ListDebtFormat,
				marker,
				strings.ToUpper(debt.Description),
				formatMoney(debt.Amount),
				debtStatus(debt),
			),
		)
	}

	sb.WriteString(manager.SpiralDelimiter)
	sb.WriteString(fmt.Sprintf(manager.ListTotalAmountFormat, formatMoney(total)))
	sb.WriteString(manager.SpiralDelimiter)

	sb.WriteString(manager.MotivationalPhrases[rand.Intn(len(manager.MotivationalPhrases))] + "\n\n")
	sb.WriteString(manager.SpiralDelimiter)

	return h.tg.SendMessageWithKeyboard(
		ctx,
		chatID,
		sb.String(),
		h.menuKeyBoard,
	)
}

func (h *Handler) getSession(ctx context.Context, userID, chatID int) (*session.Session, *manager.State, error) {
	ses, exists := h.sesMng.Get(ctx, userID)
	if !exists {
		err := h.tg.SendMessageWithKeyboard(ctx, chatID, manager.SessionLost, h.menuKeyBoard)
		if err != nil {
			h.logger.Errorf("failed to send lost session for user %d", userID)
		}

		return nil, nil, fmt.Errorf("session not found")
	}

	state, err := manager.ExtractState(ses)
	if err != nil {
		h.logger.Errorf("failed to extract state for user %d", userID)

		err = h.tg.SendMessageWithKeyboard(ctx, chatID, manager.FailedToGetState, h.menuKeyBoard)
		if err != nil {
			h.logger.Errorf("failed to send state for user %d", userID)
		}

		return nil, nil, err
	}

	return ses, state, nil
}

func debtStatus(debt *model.Debt) string {
	if debt.ReturnDate == nil {
		return manager.ReturnDateNil
	}

	days := int(time.Until(*debt.ReturnDate).Hours() / 24)
	if days < 0 {
		return fmt.Sprintf(manager.ListReturnDateExpiredFormat, -days)
	}
	return fmt.Sprintf(manager.ListReturnDateFormat, days)
}

func formatMoney(amount int64) string {
	str := fmt.Sprintf("%d", amount)
	var res []byte
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			res = append(res, '.')
		}
		res = append(res, byte(c))
	}
	return string(res)
}

func sortDebts(debts []*model.Debt) []*model.Debt {
	sort.Slice(debts, func(i, j int) bool {
		now := time.Now()

		if debts[i].ReturnDate != nil && debts[j].ReturnDate != nil {
			iOverdue := debts[i].ReturnDate.Before(now)
			jOverdue := debts[j].ReturnDate.Before(now)

			if iOverdue && jOverdue {
				return debts[i].ReturnDate.Before(*debts[j].ReturnDate)
			}

			if iOverdue {
				return false
			}
			if jOverdue {
				return true
			}

			return debts[i].ReturnDate.Before(*debts[j].ReturnDate)
		}

		if debts[i].ReturnDate == nil {
			return false
		}
		if debts[j].ReturnDate == nil {
			return true
		}
		return true
	})

	return debts
}

func (h *Handler) cleanupSession(ctx context.Context, userID int) {
	if err := h.sesMng.Delete(ctx, userID); err != nil {
		h.logger.Errorf("failed to delete session for user: %d", userID)
	}
}
