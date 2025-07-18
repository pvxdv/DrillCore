package debt

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) handleDateFlow(chatID, userID int, cmd, step, data string) error {
	switch step {
	case "year", "start":
		if data == "" {
			return h.sendWithKeyboard(chatID,
				"⏳ INITIATE TEMPORAL DRILL!\n"+
					"SELECT DESTINATION YEAR:",
				h.yearKeyboard())
		}

		return h.processYearSelection(chatID, userID, cmd, data)

	case "month":
		if data == "" {
			return h.sendWithKeyboard(chatID,
				"🌀 TEMPORAL COORDINATES PARTIAL!\n"+
					"SELECT DESTINATION MONTH:",
				h.monthKeyboard())
		}

		return h.processMonthSelection(chatID, userID, cmd, data)

	case "day":
		if data == "" {
			return h.sendWithKeyboard(chatID,
				"💢 FINAL TEMPORAL ADJUSTMENT!\n"+
					"SELECT D-DAY FOR DEBT RECLAMATION:",
				h.dateKeyboard())
		}

		return h.processDaySelection(chatID, userID, cmd, data)

	default:
		return h.sendErrorMessage(chatID,
			"🚨 UNKNOWN TEMPORAL DRILL SEQUENCE: "+step)
	}
}

func (h *Handler) processYearSelection(chatID, userID int, flow, year string) error {
	h.logger.Debugf("processing year selection: %s, for userId:%d", year, userID)
	yearInt, err := strconv.ParseInt(year, 10, 32)
	if err != nil {
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, "TEMPORAL ANOMALY DETECTED!\n"+
			"INVALID YEAR FORMAT! TRY AGAIN! 🚨")
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
				fmt.Sprintf("🌀 YEAR %s LOCKED!\n"+
					"NOW SELECT DESTINATION MONTH:", year),
				h.monthKeyboard())
		}
	}
	h.sessionMgr.Delete(userID)
	return h.sendErrorMessage(chatID, "🚨 TEMPORAL DRILL FAILURE! SESSION LOST!")
}

func (h *Handler) processMonthSelection(chatID, userID int, flow, month string) error {
	h.logger.Debugf("processing month selection: %s, for userId:%d", month, userID)
	monthInt, err := strconv.ParseInt(month, 10, 32)
	if err != nil {
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID,
			"💥 TEMPORAL ANOMALY DETECTED!\n"+
				"INVALID MONTH FORMAT! TRY AGAIN!")
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
				fmt.Sprintf("💢 MONTH %s ENGAGED!\n"+
					"NOW SET FINAL D-DAY COORDINATES:",
					time.Month(monthInt).String()),
				h.dayKeyboard(int(monthInt)))
		}
	}

	h.sessionMgr.Delete(userID)
	return h.sendErrorMessage(chatID, "TEMPORAL DRILL FAILURE! SESSION LOST! 🚨")
}

func (h *Handler) processDaySelection(chatID, userID int, flow, day string) error {
	h.logger.Debugf("processing day selection: %s, for userId:%d", day, userID)

	dayInt, err := strconv.ParseInt(day, 10, 32)
	if err != nil {
		h.sessionMgr.Delete(userID)

		return h.sendErrorMessage(chatID,
			"TEMPORAL ANOMALY DETECTED!\n"+
				"INVALID DAY FORMAT! TRY AGAIN! 🚨")
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
	}

	h.sessionMgr.Delete(userID)
	return h.sendErrorMessage(chatID, "🚨 TEMPORAL DRILL FAILURE! SESSION LOST!")
}

func (h *Handler) finishDateFlow(chatID, userID int, flowType string) error {
	if session, exists := h.sessionMgr.Get(userID); exists {
		if s, ok := session.State.(*State); ok {

			switch flowType {
			case flowAdd:
				_, err := h.storage.Save(h.ctx, s.TempDebt)
				if err != nil {
					h.logger.Debugf("failed to save debt: %v", err)

					h.sessionMgr.Delete(userID)
					return h.sendErrorMessage(chatID,
						"💥 COSMIC DEBT REGISTRY REJECTED OUR DRILL!")
				}

				h.sessionMgr.Delete(userID)

				return h.sendWithKeyboard(chatID,
					fmt.Sprintf("⚡ DEBT DRILL LAUNCH SUCCESS!\n\n"+
						"▫️ TARGET LOCKED: %s\n"+
						"▫️ SPIRAL ENERGY: %s\n"+
						"▫️ TERMINATION DATE: %s\n\n"+
						"💢 THIS MISSION HAS BEEN CARVED INTO THE BATTLE LOG!",
						strings.ToUpper(s.TempDebt.Description),
						formatMoney(s.TempDebt.Amount),
						s.TempDebt.ReturnDate.Format("02.01.2006")),
					h.debtsKeyboard())

			case flowEdit:
				s.Step = "choice"
				h.sessionMgr.Set(userID, h.ID(), s)

				return h.sendWithKeyboard(chatID,
					"🌀 TEMPORAL COORDINATES UPDATED!\n"+
						"SELECT NEXT REALITY FRAGMENT TO MODIFY:",
					h.editOptionsKeyboard())

			default:
				h.sessionMgr.Delete(userID)

				return h.sendErrorMessage(chatID,
					"🚨 UNKNOWN TEMPORAL OPERATION TYPE!")
			}
		}
	}

	return h.sendErrorMessage(chatID, "🚨 TEMPORAL DRILL FAILURE! SESSION LOST!")
}
