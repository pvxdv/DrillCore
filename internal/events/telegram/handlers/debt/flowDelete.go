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
			h.logger.Errorf("failed to get debts: %v", err)

			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "COSMIC DEBT RADAR OFFLINE! 🚨")
		}

		if len(debts) == 0 {
			h.sessionMgr.Delete(userID)
			return h.sendWithKeyboard(chatID,
				"🌟 YOUR DEBT FIELD IS CLEAR! NOTHING TO ANNIHILATE!",
				h.debtsKeyboard())
		}

		newState := &State{
			FlowType: flowDelete,
			Step:     "select",
		}
		h.sessionMgr.Set(userID, h.ID(), newState)

		return h.sendWithKeyboard(chatID,
			"💥 SELECT TARGET FOR TOTAL ANNIHILATION:",
			h.debtsListKeyboard(debts, flowDelete))

	case "select":
		debtID, err := strconv.ParseInt(strings.TrimPrefix(data, "debt_delete_select_"), 10, 64)
		if err != nil {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "INVALID TARGET IDENTIFIER! 🚨")
		}

		debt, err := h.storage.Get(h.ctx, debtID)
		if err != nil {
			h.logger.Errorf("failed to get debt %d: %v", debtID, err)
			return h.sendErrorMessage(chatID, "TARGET LOCK FAILED! 🚨")
		}

		h.logger.Debugf("get temp debt:%+v", debt)

		if debt.UserID != int64(userID) {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "💢 DRILL COLLISION DETECTED! THIS DEBT CORE BELONGS TO ANOTHER PILOT! "+
				"YOUR DRILL CANNOT PIERCE ANOTHER MAN'S SOUL! ⚔️")
		}

		state := &State{
			FlowType: flowDelete,
			TempDebt: debt,
			Step:     "confirm",
		}
		h.sessionMgr.Set(userID, h.ID(), state)

		return h.sendWithKeyboard(chatID,
			fmt.Sprintf("☠️ FINAL DRILLING SEQUENCE INITIATED!\n\n"+
				"▫️ TARGET: %s\n"+
				"▫️ DEBT LOAD: %s₽\n\n"+
				"💢 ENGAGE TOTAL ANNIHILATION?",
				debt.Description,
				formatMoney(debt.Amount)),
			h.confirmDeleteKeyboard())

	case "confirm":
		if data == cbDeleteConfirm {
			state, ok := session.State.(*State)
			if !ok || state.TempDebt == nil {
				return h.sendErrorMessage(chatID, "DRILL SEQUENCE CORRUPTED! 🚨")
			}

			if err := h.storage.Delete(h.ctx, state.TempDebt.ID); err != nil {
				h.logger.Errorf("failed to delete debt %d: %v", state.TempDebt.ID, err)
				h.sessionMgr.Delete(userID)
				return h.sendErrorMessage(chatID, "COSMIC ERASURE FAILED! 🚨")
			}

			h.sessionMgr.Delete(userID)
			return h.sendWithKeyboard(chatID,
				fmt.Sprintf("💀 TARGET DESTROYED!\n\n"+
					"▫️ %s\n"+
					"▫️ %s₽ DEBT LOAD ERASED FROM EXISTENCE!\n\n"+
					"THE DRILL PIERCED EVEN OBLIVION!",
					strings.ToUpper(state.TempDebt.Description),
					formatMoney(state.TempDebt.Amount)),
				h.debtsKeyboard())

		} else if data == cbCancel {
			h.sessionMgr.Delete(userID)
			return h.sendWithKeyboard(chatID, "❌ Удаление отменено", h.debtsKeyboard())
		}

	default:
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("🚨 UNKNOWN DRILL SEQUENCE: %s", step))
	}

	return nil
}
