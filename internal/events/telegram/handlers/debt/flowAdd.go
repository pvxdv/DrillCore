package debt

import (
	"drillCore/internal/model"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) handleAddFlow(chatID, userID int, step string, data string) error {
	session, exists := h.sessionMgr.Get(userID)
	if !exists {
		newState := &State{
			FlowType: flowAdd,
			TempDebt: &model.Debt{UserID: int64(userID)},
			Step:     stepDescription,
		}
		h.sessionMgr.Set(userID, h.ID(), newState)
		return h.sendWithKeyboard(chatID, "Введите описание долга:", h.cancelKeyboard())
	}

	state, ok := session.State.(*State)
	if !ok {
		return h.sendErrorMessage(chatID, fmt.Sprintf("unknown step: %s"+step))
	}

	switch state.Step {
	case stepDescription:
		if strings.TrimSpace(data) == "" {
			return h.sendWithKeyboard(chatID, "Описание не может быть пустым:", h.cancelKeyboard())
		}
		state.TempDebt.Description = data
		state.Step = stepAmount
		h.sessionMgr.Set(userID, h.ID(), state)
		return h.sendWithKeyboard(chatID, "Введите сумму долга:", h.cancelKeyboard())

	case stepAmount:
		amount, err := strconv.ParseInt(data, 10, 64)
		if err != nil || amount <= 0 {
			return h.sendWithKeyboard(chatID, "Неверная сумма. Введите число:", h.cancelKeyboard())
		}
		state.TempDebt.Amount = amount
		state.Step = stepYear
		now := time.Now()
		state.TempDebt.ReturnDate = &now
		h.sessionMgr.Set(userID, h.ID(), state)
		return h.sendWithKeyboard(chatID, "Установите дату возврата:", h.dateKeyboard())

	default:
		return h.sendErrorMessage(chatID, fmt.Sprintf("Неизвестный шаг в эдд флоу:%s", step))
	}
}
