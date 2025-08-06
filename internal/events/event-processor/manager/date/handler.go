package date

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	return manager.DateHandler
}

var (
	weekdays = []string{
		"MON",
		"TUE",
		"WED",
		"THU",
		"FRI",
		"SAT",
		"SUN",
	}

	monthOrder = []string{
		"🌀 JAN",
		"🌀 FEB",
		"🌀 MAR",
		"🌀 APR",
		"🌀 MAY",
		"🌀 JUN",
		"🌀 JUL",
		"🌀 AUG",
		"🌀 SEP",
		"🌀 OCT",
		"🌀 NOV",
		"🌀 DEC",
	}

	monthButtonMap = map[string]time.Month{
		"🌀 JAN": time.January,
		"🌀 FEB": time.February,
		"🌀 MAR": time.March,
		"🌀 APR": time.April,
		"🌀 MAY": time.May,
		"🌀 JUN": time.June,
		"🌀 JUL": time.July,
		"🌀 AUG": time.August,
		"🌀 SEP": time.September,
		"🌀 OCT": time.October,
		"🌀 NOV": time.November,
		"🌀 DEC": time.December,
	}
)

func (h *Handler) Handle(ctx context.Context, e *events.Event) error {
	h.logger.Debugw(
		"handling event in ",
		"handler", manager.DateHandler,
		"event", e,
	)

	if e.Type != events.Callback {
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.InvalidEventType,
			),
		)
	}

	ses, exists := h.sesMng.Get(ctx, e.Meta.UserID)
	if !exists {
		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.SessionLost,
			),
		)
	}

	cb, err := manager.ParseCallBack(e.Text)
	if err != nil {
		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToGetCallBack,
		)
	}

	state, err := manager.ExtractState(ses)
	if err != nil {
		h.logger.Errorf(
			"failed to extract state for user %d",
			e.Meta.UserID,
		)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx, e.Meta.ChatID,
			manager.FailedToGetState,
		)
	}

	if state.NextHandler == manager.IgnoreHandler ||
		state.NextStep == manager.StepIgnore ||
		state.BackStep == manager.StepIgnore ||
		state.BackHandler == manager.IgnoreHandler {
		h.cleanupSession(ctx, e.Meta.ChatID)
		return h.tg.SendMessage(
			ctx, e.Meta.ChatID,
			manager.MsgFailedRedirection,
		)
	}

	switch cb.Step {
	case manager.StepYear:
		return h.year(ctx, e, state, ses, cb)

	case manager.StepMonth:
		return h.month(ctx, e, state, ses, cb)

	case manager.StepDay:
		return h.day(ctx, e, state, ses, cb)

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

func (h *Handler) year(ctx context.Context, e *events.Event, state *manager.State, ses *session.Session, cb *manager.CallBack) error {
	kb, err := h.yearKeyboard(state.BackHandler, state.BackStep)
	if err != nil {
		h.logger.Errorf("failed to create year keyboard for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	if cb.Data == "" {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgSetYear,
			kb,
		)
	}

	yearInt, err := strconv.ParseInt(cb.Data, 10, 32)
	if err != nil {
		h.logger.Errorf("failed to parse year: %s for user: %d", cb.Data, e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidYear,
			kb,
		)
	}

	now := time.Now().UTC()
	newDate := time.Date(int(yearInt), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	if int(yearInt) < now.Year() {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.MsgDateInPast,
				formatTemporal(now, true, false, false),
				formatTemporal(newDate, true, false, false),
			),
			kb,
		)
	}

	state.TempDate = &newDate
	ses.State = state

	h.logger.Debugf("update temp date with new year: %+v:", state.TempDate)

	err = h.sesMng.Set(ctx, e.Meta.ChatID, ses)
	if err != nil {
		h.logger.Errorf("failed to set session for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(ctx, e.Meta.ChatID, manager.MsgFailedToSetSession)
	}

	kb, err = h.monthKeyboard(state.BackHandler, state.BackStep)
	if err != nil {
		h.logger.Errorf("failed to create month keyboard for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		e.Meta.ChatID,
		fmt.Sprintf(
			manager.MsgSetMonth,
			yearInt,
		),
		kb,
	)
}

func (h *Handler) month(ctx context.Context, e *events.Event, state *manager.State, ses *session.Session, cb *manager.CallBack) error {
	kb, err := h.monthKeyboard(state.BackHandler, state.BackStep)
	if err != nil {
		h.logger.Errorf("failed to create month keyboard for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	if cb.Data == "" {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.MsgSetMonth,
				state.TempDate.Year(),
			),
			kb,
		)
	}

	month, exists := monthButtonMap[cb.Data]
	if !exists {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidMonth,
			kb,
		)
	}

	now := time.Now().UTC()
	selectedYear := state.TempDate.Year()
	newDate := time.Date(selectedYear, month, 1, 0, 0, 0, 0, time.UTC)

	if selectedYear < now.Year() || (selectedYear == now.Year() && month < now.Month()) {
		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.MsgDateInPast,
				formatTemporal(now, true, true, false),
				formatTemporal(newDate, true, true, false),
			),
			kb,
		)
	}

	state.TempDate = &newDate
	ses.State = state

	h.logger.Debugf("update temp date with new month: %+v:", state.TempDate)

	err = h.sesMng.Set(ctx, e.Meta.ChatID, ses)
	if err != nil {
		h.logger.Errorf("failed to set session for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(ctx, e.Meta.ChatID, manager.MsgFailedToSetSession)
	}

	kb, err = h.dayKeyboard(selectedYear, month, state.BackHandler, state.BackStep)
	if err != nil {
		h.logger.Errorf("failed to create month keyboard for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		e.Meta.ChatID,
		fmt.Sprintf(
			manager.MsgSetDay,
			strings.ToUpper(month.String()[:3]),
		),
		kb,
	)
}

func (h *Handler) day(ctx context.Context, e *events.Event, state *manager.State, ses *session.Session, cb *manager.CallBack) error {
	kb, err := h.dateKeyboard(state.NextHandler)
	if err != nil {
		h.logger.Errorf("failed to create month keyboard for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	if cb.Data == "" {
		h.logger.Errorf("failed to get data for day step for user: %d", e.Meta.ChatID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgEmptyDay,
			kb,
		)
	}

	dayInt, err := strconv.ParseInt(cb.Data, 10, 32)
	if err != nil {
		h.logger.Errorf("failed to parse day: %s for user: %d", cb.Data, e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			manager.MsgInvalidDay,
			kb,
		)
	}

	now := time.Now().UTC()
	newDate := time.Date(state.TempDate.Year(), state.TempDate.Month(), int(dayInt), 0, 0, 0, 0, time.UTC)

	if newDate.Before(now) {
		dayKb, err := h.dayKeyboard(state.TempDate.Year(), state.TempDate.Month(), state.BackHandler, state.BackStep)
		if err != nil {
			h.logger.Errorf("failed to create month keyboard for user: %d", e.Meta.ChatID)

			h.cleanupSession(ctx, e.Meta.ChatID)

			return h.tg.SendMessage(
				ctx,
				e.Meta.ChatID,
				manager.FailedToCreateKeyboard,
			)
		}

		return h.tg.SendMessageWithKeyboard(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				manager.MsgDateInPast,
				formatTemporal(now, true, true, true),
				formatTemporal(newDate, true, true, true),
			),
			dayKb,
		)
	}

	state.TempDate = &newDate
	ses.State = state

	h.logger.Debugf("update temp debt with new day: %+v:", state.TempDate)

	err = h.sesMng.Set(ctx, e.Meta.ChatID, ses)
	if err != nil {
		h.logger.Errorf("failed to set session for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(ctx, e.Meta.ChatID, manager.MsgFailedToSetSession)
	}

	kb, err = h.redirectKeyboard(state.BackHandler, state.BackStep, state.NextHandler, state.NextStep)
	if err != nil {
		h.logger.Errorf("failed to create redirect keyboard for user: %d", e.Meta.ChatID)

		h.cleanupSession(ctx, e.Meta.ChatID)

		return h.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			manager.FailedToCreateKeyboard,
		)
	}

	return h.tg.SendMessageWithKeyboard(
		ctx,
		e.Meta.ChatID,
		fmt.Sprintf(
			manager.MsgRedirect,
			state.TempDate.Year(),
			state.TempDate.Month().String(),
			state.TempDate.Day(),
		),
		kb,
	)
}

func formatTemporal(date time.Time, showYear, showMonth, showDay bool) string {
	var parts []string

	if showMonth {
		parts = append(parts, strings.ToUpper(date.Month().String()[:3]))
	}
	if showDay {
		parts = append(parts, strconv.Itoa(date.Day()))
	}
	if showYear {
		parts = append(parts, strconv.Itoa(date.Year()))
	}

	return strings.Join(parts, " ")
}

func (h *Handler) cleanupSession(ctx context.Context, userID int) {
	if err := h.sesMng.Delete(ctx, userID); err != nil {
		h.logger.Errorf("failed to delete session for user: %d", userID)
	}
}
