package debt

import (
	"context"
	tgClient "drillCore/internal/clients/telergam"
	"drillCore/internal/events"
	"drillCore/internal/events/telegram"
	"drillCore/internal/model"
	debtStorage "drillCore/internal/storage/debt"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"sort"
	"strings"
	"time"
)

const (
	prefix = "debt_"

	cmdMenu   = "menu"
	cmdList   = "list"
	cmdCancel = "cancel"

	flowAdd    = "add"
	flowEdit   = "edit"
	flowDelete = "delete"
	flowDate   = "date"

	stepDescription = "description"
	stepAmount      = "amount"
	stepDate        = "date"
	stepYear        = "year"
	stepMonth       = "month"
	stepDay         = "day"
)

const (
	cbMenu   = prefix + cmdMenu
	cbList   = prefix + cmdList
	cbCancel = prefix + cmdCancel

	cbAddStart    = prefix + flowAdd + "_start"
	cbEditStart   = prefix + flowEdit + "_start"
	cbDeleteStart = prefix + flowDelete + "_start"
	cbDateStart   = prefix + flowDate + "_start"

	cbYear  = prefix + flowDate + "_" + stepYear + "_%d"
	cbMonth = prefix + flowDate + "_" + stepMonth + "_%d"
	cbDay   = prefix + flowDate + "_" + stepDay + "_%d"

	cbEditSelect = prefix + flowEdit + "_select_%d" // debt_edit_select_1
	cbEditDesc   = prefix + flowEdit + "_description"
	cbEditAmount = prefix + flowEdit + "_amount"

	cbEditDescProcess   = prefix + flowEdit + "_description_process"
	cbEditAmountProcess = prefix + flowEdit + "_amount_process"
	cbEditDateProcess   = prefix + flowEdit + "_date_process"

	cbEditDate   = prefix + flowEdit + "_date"
	cbFinishEdit = prefix + "edit_finish"

	cbDeleteConfirm = "debt_delete_confirm"

	flowPay      = "pay"
	cbPayConfirm = "debt_pay_confirm"

	cbPayStart = "debt_pay_start"
)

type State struct {
	FlowType string
	TempDebt *model.Debt
	Step     string
}

type Handler struct {
	tg         *tgClient.Client
	logger     *zap.SugaredLogger
	storage    debtStorage.DebtStorage
	ctx        context.Context
	sessionMgr *events.SessionManager
}

func New(ctx context.Context, tg *tgClient.Client, sm *events.SessionManager, storage debtStorage.DebtStorage, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		logger:     logger,
		storage:    storage,
		tg:         tg,
		ctx:        ctx,
		sessionMgr: sm,
	}
}

func (h *Handler) ID() string {
	return prefix
}

func (h *Handler) CanHandle(event events.Event) bool {
	if event.Type == events.Callback {
		return strings.HasPrefix(event.Text, prefix)
	}

	if event.Type == events.Message {
		meta, err := meta(event)
		if err != nil {
			return false
		}
		if session, exists := h.sessionMgr.Get(meta.UserID); exists {
			return session.HandlerID == h.ID()
		}
	}

	return false
}

func (h *Handler) Handle(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("get meta: %w", err)
	}

	h.logger.Debugf("handling event: event=%+v, Text='%s'", event, event.Text)

	if strings.HasPrefix(event.Text, cbCancel) {
		h.sessionMgr.Delete(meta.UserID)
		return h.showDebtMenu(meta.ChatID)
	}

	if session, exists := h.sessionMgr.Get(meta.UserID); exists {
		if s, ok := session.State.(*State); ok {
			return h.handleFlow(event, meta, s)
		}
	}

	cmdRaw := strings.TrimPrefix(event.Text, prefix)
	if cmdRaw == "" {
		h.sessionMgr.Delete(meta.UserID)
		h.logger.Debugf("no command")
		return h.showDebtMenu(meta.ChatID)
	}

	cmdParts := strings.Split(cmdRaw, "_")
	if len(cmdParts) == 0 {
		h.sessionMgr.Delete(meta.UserID)
		return h.showDebtMenu(meta.ChatID)
	}

	h.logger.Debugf("handling command: %v", cmdParts)
	switch cmdParts[0] {
	case cmdMenu:
		return h.showDebtMenu(meta.ChatID)
	case cmdList:
		return h.listDebts(meta.ChatID, meta.UserID)
	case flowAdd, flowEdit, flowDelete, flowDate, flowPay:
		return h.handleFlow(event, meta, nil)
	default:
		return h.sendErrorMessage(meta.ChatID, fmt.Sprintf("unknown command: %s"+cmdRaw))
	}
}

func (h *Handler) handleFlow(event events.Event, meta telegram.Meta, state *State) error {
	var cmd string
	var step string
	var data string

	var cmdOrgin string
	if state != nil {
		cmd = state.FlowType
		step = state.Step
		if step == stepYear || step == stepMonth || step == stepDay {
			cmdOrgin = cmd
			cmd = flowDate
		}

		if step == "choice" {
			step = ""
		}

		h.logger.Debugf("handling flow: event:%v, flow:%s, step:%s", event.Text, state.FlowType, step)
	}

	cmdRaw := strings.TrimPrefix(event.Text, prefix)
	if cmdRaw == "" {
		h.sessionMgr.Delete(meta.UserID)
		h.logger.Debugf("no command")
		return h.showDebtMenu(meta.ChatID)
	}

	cmdParts := strings.Split(cmdRaw, "_")
	if len(cmdParts) == 0 {
		h.sessionMgr.Delete(meta.UserID)
		return h.showDebtMenu(meta.ChatID)
	}

	if len(cmdParts) > 0 && state == nil {
		cmd = cmdParts[0]
		cmdOrgin = cmd
	}

	if len(cmdParts) > 1 {
		if step == "" {
			step = cmdParts[1]
		}
	}

	if len(cmdParts) > 2 {
		data = cmdParts[2]
	}

	h.logger.Debugf("checking command: %v, step:%v, data:%v", cmd, step, data)
	switch cmd {
	case flowAdd:
		return h.handleAddFlow(meta.ChatID, meta.UserID, step, event.Text)
	case flowDate:
		return h.handleDateFlow(meta.ChatID, meta.UserID, cmdOrgin, step, data)
	case flowEdit:
		return h.handleEditFlow(meta.ChatID, meta.UserID, step, event.Text)
	case flowDelete:
		return h.handleDeleteFlow(meta.ChatID, meta.UserID, step, event.Text)
	case flowPay:
		return h.handlePayFlow(meta.ChatID, meta.UserID, step, event.Text)

	default:
		return h.showDebtMenu(meta.ChatID)
	}
}

func (h *Handler) showDebtMenu(chatID int) error {
	return h.sendWithKeyboard(chatID, "Выберите действие:", h.debtsKeyboard())
}

func (h *Handler) listDebts(chatID, userID int) error {
	debts, err := h.storage.Debts(h.ctx, int64(userID))
	if err != nil {
		h.sessionMgr.Delete(userID)
		h.logger.Errorf("SPIRAL POWER FAILURE: %v", err)
		return h.sendErrorMessage(chatID, "🚨 DEBT BRILL SYSTEM OFFLINE! 🚨")
	}

	if len(debts) == 0 {
		var sb strings.Builder
		sb.WriteString("─────────🌀─────────\n\n")
		sb.WriteString(noDebtsPhrases[rand.Intn(len(noDebtsPhrases))] + "\n\n")
		sb.WriteString("─────────🌀─────────\n\n")

		h.sessionMgr.Delete(userID)
		return h.sendWithKeyboard(
			chatID,
			sb.String(),
			h.debtsKeyboard(),
		)
	}

	// Сортируем: сначала просроченные (враги ближе всего), затем по дате
	sort.Slice(debts, func(i, j int) bool {
		now := time.Now()
		iOverdue := debts[i].ReturnDate != nil && debts[i].ReturnDate.Before(now)
		jOverdue := debts[j].ReturnDate != nil && debts[j].ReturnDate.Before(now)

		if iOverdue && !jOverdue {
			return true
		}
		if !iOverdue && jOverdue {
			return false
		}
		if debts[i].ReturnDate == nil {
			return false
		}
		if debts[j].ReturnDate == nil {
			return true
		}
		return debts[i].ReturnDate.Before(*debts[j].ReturnDate)
	})

	total := int64(0)
	for _, d := range debts {
		total += d.Amount
	}

	var sb strings.Builder
	sb.WriteString("─────────🌀─────────\n\n")
	sb.WriteString(debtTitles[rand.Intn(len(debtTitles))] + "\n\n")
	sb.WriteString("─────────🌀─────────\n\n") // Спиральный разделитель

	for i, debt := range debts {
		// Чередующиеся маркеры в стиле Team Gurren
		marker := "💢" // Нечетные
		if i%2 == 0 {
			marker = "🌀" // Четные
		}
		if debt.ReturnDate != nil && debt.ReturnDate.Before(time.Now()) {
			marker = "☠️" // Просрочки - череп
		}

		sb.WriteString(fmt.Sprintf(
			"%s %s\n"+
				"   💥 SPIRAL COST: %s₽\n"+ // Изменено на SPIRAL COST
				"   %s\n\n",
			marker,
			strings.ToUpper(debt.Description),
			formatMoney(debt.Amount),
			getDebtStatus(debt)))
	}

	sb.WriteString("─────────🌀─────────\n\n")
	sb.WriteString(fmt.Sprintf("💢 TOTAL SPIRAL POWER NEEDED: %s₽\n\n", formatMoney(total)))
	sb.WriteString("─────────🌀─────────\n\n")

	sb.WriteString(motivationalPhrases[rand.Intn(len(motivationalPhrases))] + "\n\n")
	sb.WriteString("─────────🌀─────────")

	return h.sendWithKeyboard(chatID, sb.String(), h.debtsKeyboard())
}

func getDebtStatus(debt *model.Debt) string {
	if debt.ReturnDate == nil {
		return "🌌 UNLIMITED BATTLEFIELD"
	}

	days := int(time.Until(*debt.ReturnDate).Hours() / 24)
	if days < 0 {
		return fmt.Sprintf("🚨 ANTI-SPIRAL THREAT (%d DAYS)", -days)
	}
	return fmt.Sprintf("⏳ GIGA DRILL CHARGE: %d DAYS", days)
}

// Функция форматирования денег (оставляем без изменений)
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

func meta(event events.Event) (telegram.Meta, error) {
	res, ok := event.Meta.(telegram.Meta)
	if !ok {
		return telegram.Meta{}, fmt.Errorf("failed to process meta: %w", telegram.ErrUnknownMetaType)
	}
	return res, nil
}

func (h *Handler) sendMessage(chatID int, text string) error {
	return h.tg.SendMessage(chatID, text)
}

func (h *Handler) sendWithKeyboard(chatID int, text string, keyboard tgClient.ReplyMarkup) error {
	return h.tg.SendMessageWithKeyboard(chatID, text, keyboard)
}

func (h *Handler) sendErrorMessage(chatID int, message string) error {
	return h.sendWithKeyboard(chatID, "⚠️ "+message, h.debtsKeyboard())
}
