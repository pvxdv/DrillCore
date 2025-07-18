package debt

import (
	"fmt"
	"strconv"
	"time"
)

func (h *Handler) handleDateFlow(chatID, userID int, cmd, step, data string) error {
	switch step {
	case "year", "start":
		if data == "" {
			return h.sendWithKeyboard(chatID, "выберете год", h.yearKeyboard())
		}

		return h.processYearSelection(chatID, userID, cmd, data)

	case "month":
		if data == "" {
			return h.sendWithKeyboard(chatID, "выберете месяц", h.monthKeyboard())
		}

		return h.processMonthSelection(chatID, userID, cmd, data)

	case "day":
		if data == "" {
			return h.sendWithKeyboard(chatID, "выберете день", h.dateKeyboard())
		}

		return h.processDaySelection(chatID, userID, cmd, data)

	default:
		return h.sendErrorMessage(chatID, fmt.Sprintf("unknown step: %s"+step))
	}
}

func (h *Handler) processYearSelection(chatID, userID int, flow, year string) error {
	h.logger.Debugf("processing year selection: %s, for userId:%d", year, userID)
	yearInt, err := strconv.ParseInt(year, 10, 32)
	if err != nil {
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("invalid year: %v", err))
	}

	now := time.Now()

	newDate := time.Date(
		int(yearInt),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	if session, exists := h.sessionMgr.Get(userID); exists {
		if s, ok := session.State.(*State); ok {
			s.TempDebt.ReturnDate = &newDate
			s.FlowType = flow
			s.Step = stepMonth

			h.logger.Debugf("update temp debt with new year: %+v:", s.TempDebt)

			h.sessionMgr.Set(userID, h.ID(), s)

			return h.sendWithKeyboard(chatID,
				fmt.Sprintf("Выберите месяц для %s года:", year),
				h.monthKeyboard())
		}
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("invalid step, session not found: %v", err))
	}
	h.sessionMgr.Delete(userID)
	return h.sendErrorMessage(chatID, fmt.Sprintf("invalid step, session not found: %v", err))
}

func (h *Handler) processMonthSelection(chatID, userID int, flow, month string) error {
	h.logger.Debugf("processing month selection: %s, for userId:%d", month, userID)
	monthInt, err := strconv.ParseInt(month, 10, 32)
	if err != nil {
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("invalid month: %v", err))
	}

	if session, exists := h.sessionMgr.Get(userID); exists {
		if s, ok := session.State.(*State); ok {
			old := s.TempDebt.ReturnDate

			now := time.Now()

			newDate := time.Date(
				old.Year(),
				time.Month(int(monthInt)),
				now.Day(),
				0, 0, 0, 0,
				time.UTC,
			)

			s.TempDebt.ReturnDate = &newDate
			s.FlowType = flow
			s.Step = stepDay

			h.logger.Debugf("update temp debt with new month: %+v:", s.TempDebt)

			h.sessionMgr.Set(userID, h.ID(), s)

			return h.sendWithKeyboard(chatID,
				fmt.Sprintf("Выберите день для %d года и %d месяца", old.Year(), newDate.Month()),
				h.dayKeyboard(int(monthInt)))
		}
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("invalid step, session not found: %v", err))
	}
	h.sessionMgr.Delete(userID)
	return h.sendErrorMessage(chatID, fmt.Sprintf("invalid step, session not found: %v", err))
}

func (h *Handler) processDaySelection(chatID, userID int, flow, day string) error {
	h.logger.Debugf("processing day selection: %s, for userId:%d", day, userID)
	dayInt, err := strconv.ParseInt(day, 10, 32)
	if err != nil {
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("invalid day: %v", err))
	}

	if session, exists := h.sessionMgr.Get(userID); exists {
		if s, ok := session.State.(*State); ok {
			old := s.TempDebt.ReturnDate

			newDate := time.Date(
				old.Year(),
				old.Month(),
				int(dayInt),
				0, 0, 0, 0,
				time.UTC,
			)

			s.TempDebt.ReturnDate = &newDate
			s.FlowType = flow
			s.Step = stepDay

			h.logger.Debugf("update temp debt with new day: %+v:", s.TempDebt)

			h.sessionMgr.Set(userID, h.ID(), s)

			return h.finishDateFlow(chatID, userID, s.FlowType)
		}
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("invalid step, session not found: %v", err))
	}
	h.sessionMgr.Delete(userID)
	return h.sendErrorMessage(chatID, fmt.Sprintf("invalid step, session not found: %v", err))
}

func (h *Handler) finishDateFlow(chatID, userID int, flowType string) error {
	var err error
	var msg string

	if session, exists := h.sessionMgr.Get(userID); exists {
		if s, ok := session.State.(*State); ok {

			switch flowType {
			case flowAdd:
				_, err = h.storage.Save(h.ctx, s.TempDebt)
				msg = "✅ Долг успешно добавлен"
			case flowEdit:
				s.Step = "choice"
				h.sessionMgr.Set(userID, h.ID(), s)
				return h.sendWithKeyboard(chatID, "дата обновлена. Что редактируем дальше?", h.editOptionsKeyboard())
			default:
				h.sessionMgr.Delete(userID)
				return h.sendErrorMessage(chatID, "Неизвестный тип операции")
			}

			h.sessionMgr.Delete(userID)

			if err != nil {
				h.logger.Errorf("Ошибка сохранения долга: %v", err)
				return h.sendErrorMessage(chatID, "Ошибка сохранения")
			}

			return h.sendWithKeyboard(chatID, msg, h.debtsKeyboard())
		}

		return h.sendErrorMessage(chatID, fmt.Sprintf("invalid step, session not found: %v", err))
	}
	h.sessionMgr.Delete(userID)
	return h.sendErrorMessage(chatID, fmt.Sprintf("invalid step, session not found: %v", err))
}
