package debt

import (
	"fmt"
	"strconv"
	"strings"
)

func (h *Handler) handleDeleteFlow(chatID, userID int, step string, data string) error {
	h.logger.Debugf("handle delete flow step: %s, data:%s", step, data)

	session, exists := h.sessionMgr.Get(userID)
	if !exists && step != "start" {
		return h.handleDeleteFlow(chatID, userID, "start", "")
	}

	switch step {
	case "start":
		debts, err := h.storage.Debts(h.ctx, int64(userID))
		if err != nil {
			h.logger.Errorf("Ошибка получения списка долгов: %v", err)
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "Ошибка получения списка долгов")
		}

		if len(debts) == 0 {
			return h.sendWithKeyboard(chatID, "У вас нет долгов для удаления", h.debtsKeyboard())
		}

		newState := &State{
			FlowType: flowDelete,
			Step:     "select",
		}
		h.sessionMgr.Set(userID, h.ID(), newState)

		return h.sendWithKeyboard(chatID,
			"Выберите долг для удаления:",
			h.debtsListKeyboard(debts, flowDelete))

	case "select":
		debtID, err := strconv.ParseInt(strings.TrimPrefix(data, "debt_delete_select_"), 10, 64)
		if err != nil {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "Неверный ID долга")
		}

		debt, err := h.storage.Get(h.ctx, debtID)
		if err != nil {
			h.logger.Errorf("Ошибка получения долга %d: %v", debtID, err)
			return h.sendErrorMessage(chatID, "Ошибка получения долга")
		}

		h.logger.Debugf("get temp debt:%+v", debt)

		if debt.UserID != int64(userID) {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "Это не ваш долг")
		}

		state := &State{
			FlowType: flowDelete,
			TempDebt: debt,
			Step:     "confirm",
		}
		h.sessionMgr.Set(userID, h.ID(), state)

		// Форматируем дату с проверкой на nil
		var dateStr string
		if debt.ReturnDate != nil {
			dateStr = debt.ReturnDate.Format("02.01.2006")
		} else {
			dateStr = "не указана"
		}

		confirmMsg := fmt.Sprintf(
			"Вы действительно хотите удалить этот долг?\n\n"+
				"🔹 Описание: %s\n"+
				"🔹 Сумма: %d₽\n"+
				"🔹 Дата возврата: %s\n\n"+
				"⚠️ Это действие необратимо!",
			debt.Description,
			debt.Amount,
			dateStr,
		)

		return h.sendWithKeyboard(chatID, confirmMsg, h.confirmDeleteKeyboard())

	case "confirm":
		if data == string(cbDeleteConfirm) {
			state, ok := session.State.(*State)
			if !ok || state.TempDebt == nil {
				return h.sendErrorMessage(chatID, "Ошибка: данные сессии утеряны")
			}

			if err := h.storage.Delete(h.ctx, state.TempDebt.ID); err != nil {
				h.logger.Errorf("Ошибка удаления долга %d: %v", state.TempDebt.ID, err)
				h.sessionMgr.Delete(userID)
				return h.sendErrorMessage(chatID, "Ошибка удаления долга")
			}

			h.sessionMgr.Delete(userID)
			return h.sendWithKeyboard(chatID, "✅ Долг успешно удален", h.debtsKeyboard())
		} else if data == string(cbCancel) {
			h.sessionMgr.Delete(userID)
			return h.sendWithKeyboard(chatID, "❌ Удаление отменено", h.debtsKeyboard())
		}

	default:
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("Неизвестный шаг в процессе удаления: %s", step))
	}

	return nil
}
