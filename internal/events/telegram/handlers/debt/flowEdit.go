package debt

import (
	"drillCore/internal/model"
	"fmt"
	"strconv"
	"strings"
)

func (h *Handler) handleEditFlow(chatID, userID int, step string, data string) error {
	h.logger.Debugf("handle edit flow step: %s, data:%s", step, data)
	session, exists := h.sessionMgr.Get(userID)
	if !exists {
		newState := &State{
			FlowType: flowEdit,
			TempDebt: &model.Debt{UserID: int64(userID)},
			Step:     "select",
		}
		h.sessionMgr.Set(userID, h.ID(), newState)

		debts, err := h.storage.Debts(h.ctx, int64(userID))
		if err != nil {
			h.sessionMgr.Delete(userID)
			h.logger.Errorf("Ошибка получения списка долгов: %v", err)
			return h.sendErrorMessage(chatID, "Ошибка получения списка долгов")
		}

		if len(debts) == 0 {
			return h.sendWithKeyboard(chatID, "У вас нет долгов для редактирования", h.debtsKeyboard())
		}

		// Отправляем список долгов с кнопками для выбора
		return h.sendWithKeyboard(chatID,
			"Выберите долг для редактирования:",
			h.debtsListKeyboard(debts, flowEdit))
	}

	state, ok := session.State.(*State)
	if !ok {
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("unknown step: %s"+step))
	}

	h.logger.Debugf("get temp debt:%+v", state.TempDebt)

	switch step {
	case "select":
		h.logger.Debugf("handle edit flow step: select, data:%s", data)
		debtID, err := strconv.ParseInt(strings.TrimPrefix(data, "debt_edit_select_"), 10, 64)
		if err != nil {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, fmt.Sprintf("Неверный ID долга:%v", err))
		}

		debt, err := h.storage.Get(h.ctx, debtID)
		if err != nil {
			h.sessionMgr.Delete(userID)
			h.logger.Errorf("Ошибка получения долга %d: %v", debtID, err)
			return h.sendErrorMessage(chatID, "Ошибка получения долга")
		}

		if debt.UserID != int64(userID) {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "Это не ваш долг")
		}

		state.TempDebt = debt
		state.Step = ""
		h.sessionMgr.Set(userID, h.ID(), state)

		return h.sendWithKeyboard(chatID, "Что хотите изменить?", h.editOptionsKeyboard())
	case "choice":
		h.logger.Debugf("handle edit flow step: choice, data:%s", data)
		return h.sendWithKeyboard(chatID, "Что хотите изменить?", h.editOptionsKeyboard())
	case stepDescription:
		h.logger.Debugf("handle edit flow step: description, data:%s", data)
		if data == "debt_edit_description_process" {
			state.Step = stepDescription
			h.sessionMgr.Set(userID, h.ID(), state)
			return h.sendWithKeyboard(chatID, "Введите новое описание:", h.cancelKeyboard())
		}

		if strings.TrimSpace(data) == "" {
			return h.sendWithKeyboard(chatID, "Описание не может быть пустым:", h.cancelKeyboard())
		}

		state.TempDebt.Description = data
		state.Step = ""
		h.sessionMgr.Set(userID, h.ID(), state)

		h.logger.Debugf("update temp debt with description: %+v:", state.TempDebt)

		return h.sendWithKeyboard(chatID, "Описание обновлено. Что редактируем дальше?", h.editOptionsKeyboard())

	case stepAmount:
		h.logger.Debugf("handle edit flow step: amount, data:%s", data)
		if data == "debt_edit_amount_process" {
			state.Step = stepAmount
			h.sessionMgr.Set(userID, h.ID(), state)
			return h.sendWithKeyboard(chatID, "Введите новую сумму:", h.cancelKeyboard())
		}

		if strings.TrimSpace(data) == "" {
			return h.sendWithKeyboard(chatID, "сумма не может быть пустой:", h.cancelKeyboard())
		}

		amount, err := strconv.ParseInt(data, 10, 64)
		if err != nil || amount <= 0 {
			return h.sendWithKeyboard(chatID, "Неверная сумма. Введите положительное число:", h.cancelKeyboard())
		}

		state.TempDebt.Amount = amount
		state.Step = ""
		h.sessionMgr.Set(userID, h.ID(), state)

		h.logger.Debugf("update temp debt with amount: %+v:", state.TempDebt)

		return h.sendWithKeyboard(chatID, "сумма обновлена. Что редактируем дальше?", h.editOptionsKeyboard())
	case stepDate:
		h.logger.Debugf("handle edit flow step: date, data:%s", data)
		if data == "debt_edit_date_process" {
			state.Step = stepYear
			h.sessionMgr.Set(userID, h.ID(), state)
			return h.sendWithKeyboard(chatID, "Установите дату возврата:", h.dateKeyboard())
		}

		if data == "" {
			state.Step = stepDate
			h.sessionMgr.Set(userID, h.ID(), state)
			return h.sendWithKeyboard(chatID, "Установите новую дату:", h.dateKeyboard())
		}
	case "finish":
		h.logger.Debugf("handle edit flow step: finish, data:%s", data)

		h.logger.Debugf("try to save:%+v", state.TempDebt)

		err := h.storage.Update(h.ctx, state.TempDebt)
		h.sessionMgr.Delete(userID)

		if err != nil {
			h.logger.Errorf("Ошибка обновления долга: %v", err)
			return h.sendErrorMessage(chatID, "Ошибка сохранения")
		}

		return h.sendWithKeyboard(chatID, "✅ Долг успешно обновлен", h.debtsKeyboard())
	default:
		return h.sendErrorMessage(chatID, fmt.Sprintf("Неизвестный шаг в эдд флоу:%s", step))
	}

	return h.sendErrorMessage(chatID, fmt.Sprintf("Неизвестный шаг в эдд флоу:%s", step))
}
