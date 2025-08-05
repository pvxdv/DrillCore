package manager

import (
	"context"
	"fmt"

	"drillCore/internal/bot"
	"drillCore/internal/events"
	"drillCore/internal/session"

	"go.uber.org/zap"
)

type Handler interface {
	Type() TypeHandler
	Handle(ctx context.Context, e *events.Event) error
}

type SessionManager interface {
	Get(ctx context.Context, userID int) (*session.Session, bool)
	Set(ctx context.Context, userID int, s *session.Session) error
	Delete(ctx context.Context, userID int) error
}

type Manager struct {
	tg       *bot.Client
	sesMng   SessionManager
	logger   *zap.SugaredLogger
	handlers map[TypeHandler]*Handler
}

func New(tg *bot.Client, sm SessionManager, logger *zap.SugaredLogger, handlers ...Handler) *Manager {
	p := &Manager{
		tg:       tg,
		logger:   logger,
		sesMng:   sm,
		handlers: registeredHandlers(handlers...),
	}

	return p
}

func (m *Manager) HandleEvent(ctx context.Context, e *events.Event) error {
	m.logger.Debugf("handle event: %+v", e)

	switch e.Type {
	case events.Message:
		return m.routeUserInput(ctx, e)
	case events.Callback:
		return m.routeCallBack(ctx, e)
	default:
		m.logger.Errorf("unknown event type %v", e.Type)

		return m.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			InvalidEventType,
		)
	}
}

func (m *Manager) routeCallBack(ctx context.Context, e *events.Event) error {
	m.logger.Debugf("route call back event: %+v", e)

	cb, err := ParseCallBack(e.Text)
	if err != nil {
		m.logger.Errorf("failed to parse callback event: %v", err)
		return m.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			FailedToGetCallBack,
		)
	}

	m.logger.Debugf("parse call back: %+v", cb)

	if cb.Handler == IgnoreHandler {
		return nil
	}

	h, ok := m.handler(cb.Handler)
	if !ok {
		m.logger.Errorf(
			"failed to find handler: %v, for user %d",
			cb.Handler,
			e.Meta.UserID,
		)

		return m.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			HandlerNotFound,
		)
	}

	m.logger.Debugf("found call back handler for user %d, h:%+v", e.Meta.ChatID, h)

	return h.Handle(ctx, e)
}

func (m *Manager) routeUserInput(ctx context.Context, e *events.Event) error {
	m.logger.Debugf("route user input: %v", e.Text)

	if _, isCmd := ParseCommand(e.Text); isCmd {
		h, ok := m.handler(CMDHandler)
		if !ok {
			m.logger.Errorf(
				"failed to find command handler, for user %d",
				e.Meta.UserID,
			)

			return m.tg.SendMessage(
				ctx,
				e.Meta.ChatID,
				fmt.Sprintf(
					InvalidStep,
					"commandHandler",
				),
			)
		}

		return h.Handle(ctx, e)
	}

	ses, exists := m.sesMng.Get(ctx, e.Meta.UserID)
	if !exists {
		m.logger.Errorf(
			"failed to find session for user %d",
			e.Meta.UserID,
		)

		return m.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			ButtonOnlyMode,
		)
	}

	state, err := ExtractState(ses)
	if err != nil {
		return err
	}

	h, ok := m.handler(state.Handler)
	if !ok {
		return m.tg.SendMessage(
			ctx,
			e.Meta.ChatID,
			fmt.Sprintf(
				InvalidStep,
				state.Step,
			),
		)
	}

	return h.Handle(ctx, e)
}

func registeredHandlers(handlers ...Handler) map[TypeHandler]*Handler {
	m := make(map[TypeHandler]*Handler, len(handlers))

	for _, h := range handlers {
		m[h.Type()] = &h
	}

	return m
}

func (m *Manager) handler(t TypeHandler) (Handler, bool) {
	h, ok := m.handlers[t]
	if !ok {
		return nil, false
	}

	return *h, ok
}
