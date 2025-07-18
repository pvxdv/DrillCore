package debt

import (
	"fmt"
	"strconv"
	"strings"
)

func (h *Handler) handlePayFlow(chatID, userID int, step string, data string) error {
	h.logger.Debugf("handle pay flow step: %s, data:%s", step, data)

	session, exists := h.sessionMgr.Get(userID)
	if !exists && step != "start" {
		return h.handlePayFlow(chatID, userID, "start", "")
	}

	switch step {
	case "start":
		debts, err := h.storage.Debts(h.ctx, int64(userID))
		if err != nil {
			h.logger.Errorf("Ошибка получения списка долгов: %v", err)
			return h.sendErrorMessage(chatID, "Ошибка получения списка долгов")
		}

		if len(debts) == 0 {
			return h.sendWithKeyboard(chatID, "У вас нет активных долгов для оплаты", h.debtsKeyboard())
		}

		newState := &State{
			FlowType: flowPay,
			Step:     "select",
		}
		h.sessionMgr.Set(userID, h.ID(), newState)

		return h.sendWithKeyboard(chatID,
			"Выберите долг для оплаты:",
			h.debtsListKeyboard(debts, flowPay))

	case "select":
		debtID, err := strconv.ParseInt(strings.TrimPrefix(data, "debt_pay_select_"), 10, 64)
		if err != nil {
			return h.sendErrorMessage(chatID, "Неверный ID долга")
		}

		debt, err := h.storage.Get(h.ctx, debtID)
		if err != nil {
			h.logger.Errorf("Ошибка получения долга %d: %v", debtID, err)
			return h.sendErrorMessage(chatID, "Ошибка получения долга")
		}

		if debt.UserID != int64(userID) {
			return h.sendErrorMessage(chatID, "Это не ваш долг")
		}

		state := &State{
			FlowType: flowPay,
			TempDebt: debt,
			Step:     "amount",
		}
		h.sessionMgr.Set(userID, h.ID(), state)

		return h.sendWithKeyboard(chatID,
			fmt.Sprintf("Текущий долг: %s\nСумма: %s₽\n\nВведите сумму оплаты:",
				debt.Description,
				formatMoney(debt.Amount)),
			h.cancelKeyboard())

	case "amount":
		state, ok := session.State.(*State)
		if !ok || state.TempDebt == nil {
			return h.sendErrorMessage(chatID, "Ошибка: данные сессии утеряны")
		}

		// Парсим сумму из data (event.Text)
		payment, err := strconv.ParseInt(data, 10, 64)
		if err != nil || payment <= 0 {
			return h.sendWithKeyboard(chatID,
				"Неверная сумма. Введите положительное число:",
				h.cancelKeyboard())
		}

		if payment > state.TempDebt.Amount {
			return h.sendWithKeyboard(chatID,
				fmt.Sprintf("Сумма превышает долг. Максимально можно внести: %s₽",
					formatMoney(state.TempDebt.Amount)),
				h.cancelKeyboard())
		}

		// Сохраняем сумму в Step (временное хранилище)
		state.Step = "confirm"
		h.sessionMgr.Set(userID, h.ID(), state)

		newAmount := state.TempDebt.Amount - payment
		confirmMsg := fmt.Sprintf(
			"Подтвердите оплату:\n\n"+
				"▪️ Долг: %s\n"+
				"▪️ Текущая сумма: %s₽\n"+
				"▪️ Оплата: %s₽\n"+
				"▪️ Остаток: %s₽\n\n"+
				"Подтверждаете?",
			state.TempDebt.Description,
			formatMoney(state.TempDebt.Amount),
			formatMoney(payment),
			formatMoney(newAmount),
		)

		return h.sendWithKeyboard(chatID, confirmMsg, h.confirmPayKeyboard(payment))

	case "confirm":
		state, ok := session.State.(*State)
		if !ok || state.TempDebt == nil {
			return h.sendErrorMessage(chatID, "Ошибка: данные сессии утеряны")
		}

		payment, err := strconv.ParseInt(strings.TrimPrefix(data, "debt_pay_confirm_"), 10, 64)
		if err != nil {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "Ошибка: неверная сумма оплаты")
		}

		newAmount := state.TempDebt.Amount - payment

		if newAmount <= 0 {
			// Полное погашение - удаляем долг
			if err := h.storage.Delete(h.ctx, state.TempDebt.ID); err != nil {
				h.logger.Errorf("Ошибка удаления долга %d: %v", state.TempDebt.ID, err)
				return h.sendErrorMessage(chatID, "Ошибка при погашении долга")
			}
			h.sessionMgr.Delete(userID)
			return h.sendWithKeyboard(chatID,
				fmt.Sprintf("✅ Долг \"%s\" полностью погашен", state.TempDebt.Description),
				h.debtsKeyboard())
		} else {
			// Частичное погашение - обновляем сумму
			state.TempDebt.Amount = newAmount
			if err := h.storage.Update(h.ctx, state.TempDebt); err != nil {
				h.logger.Errorf("Ошибка обновления долга %d: %v", state.TempDebt.ID, err)
				return h.sendErrorMessage(chatID, "Ошибка при обновлении долга")
			}
			h.sessionMgr.Delete(userID)
			return h.sendWithKeyboard(chatID,
				fmt.Sprintf("✅ Оплата проведена. Остаток по долгу \"%s\": %s₽",
					state.TempDebt.Description,
					formatMoney(newAmount)),
				h.debtsKeyboard())
		}

	default:
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("Неизвестный шаг в процессе оплаты: %s", step))
	}
}
