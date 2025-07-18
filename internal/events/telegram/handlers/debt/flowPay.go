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
			h.logger.Errorf("failed to get debts: %v", err)
			return h.sendErrorMessage(chatID, "DEBT SCAN FAILED! DEBT DRILL SYSTEM OFFLINE! 🚨")
		}

		if len(debts) == 0 {
			return h.sendWithKeyboard(chatID, "🌟 YOUR DEBT FIELD IS PRISTINE! NO TARGETS FOR PAYMENT DRILL!", h.debtsKeyboard())
		}

		newState := &State{
			FlowType: flowPay,
			Step:     "select",
		}
		h.sessionMgr.Set(userID, h.ID(), newState)

		return h.sendWithKeyboard(chatID, "💢 DEPLOY PAYMENT DRILL! SELECT YOUR TARGET:", h.debtsListKeyboard(debts, flowPay))

	case "select":
		debtID, err := strconv.ParseInt(strings.TrimPrefix(data, "debt_pay_select_"), 10, 64)
		if err != nil {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "INVALID TARGET LOCK! DEBT ID CORRUPTED! 🚨")
		}

		debt, err := h.storage.Get(h.ctx, debtID)
		if err != nil {
			h.sessionMgr.Delete(userID)
			h.logger.Errorf("failed to det debt %d: %v", debtID, err)
			return h.sendErrorMessage(chatID, "TARGET ACQUISITION FAILED! SPIRAL SIGNAL LOST! 🚨")
		}

		if debt.UserID != int64(userID) {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "💢 DRILL COLLISION DETECTED! THIS DEBT CORE BELONGS TO ANOTHER PILOT! "+
				"YOUR DRILL CANNOT PIERCE ANOTHER MAN'S SOUL! ⚔️")
		}

		state := &State{
			FlowType: flowPay,
			TempDebt: debt,
			Step:     "amount",
		}
		h.sessionMgr.Set(userID, h.ID(), state)

		return h.sendWithKeyboard(chatID,
			fmt.Sprintf(`🌀 TARGET LOCK ESTABLISHED!

▫️ TARGET DESIGNATION: %s
▫️ SPIRAL DEBT LOAD: %s₽

💢 INPUT PAYMENT DRILL ENERGY LEVEL:`,
				debt.Description,
				formatMoney(debt.Amount)),
			h.cancelKeyboard())

	case "amount":
		state, ok := session.State.(*State)
		if !ok || state.TempDebt == nil {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "SPIRAL MEMORY CORRUPTION! SEQUENCE LOST! 🚨")
		}

		payment, err := strconv.ParseInt(data, 10, 64)
		if err != nil || payment <= 0 {
			return h.sendWithKeyboard(chatID, "🚨 INVALID DRILL POWER! ENTER POSITIVE NUMBER:", h.cancelKeyboard())
		}

		if payment > state.TempDebt.Amount {
			return h.sendWithKeyboard(chatID,
				fmt.Sprintf("💥 DRILL POWER TOO STRONG! MAX: %s₽",
					formatMoney(state.TempDebt.Amount)),
				h.cancelKeyboard())
		}

		state.Step = "confirm"
		h.sessionMgr.Set(userID, h.ID(), state)

		newAmount := state.TempDebt.Amount - payment
		confirmMsg := fmt.Sprintf(
			`⚡ PAYMENT DRILL FINAL LOCK!

▫️ TARGET: %s
▫️ INITIAL DEBT LOAD: %s₽
▫️ DRILL ENERGY: %s₽
▫️ REMAINING DEBT: %s₽

💢 INITIATE SPIRAL PAYMENT SEQUENCE?`,
			state.TempDebt.Description,
			formatMoney(state.TempDebt.Amount),
			formatMoney(payment),
			formatMoney(newAmount),
		)

		return h.sendWithKeyboard(chatID, confirmMsg, h.confirmPayKeyboard(payment))

	case "confirm":
		state, ok := session.State.(*State)
		if !ok || state.TempDebt == nil {
			return h.sendErrorMessage(chatID, "SPIRAL MEMORY CORRUPTION! SEQUENCE LOST! 🚨")
		}

		payment, err := strconv.ParseInt(strings.TrimPrefix(data, "debt_pay_confirm_"), 10, 64)
		if err != nil {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "ENERGY SIGNAL DISTORTED! PAYMENT ABORTED! 🚨 ")
		}

		newAmount := state.TempDebt.Amount - payment

		if newAmount <= 0 {
			if err := h.storage.Delete(h.ctx, state.TempDebt.ID); err != nil {
				h.logger.Errorf("failed to delete debt %d: %v", state.TempDebt.ID, err)
				return h.sendErrorMessage(chatID, "DEBT ANNIHILATION FAILURE! SPIRAL COLLAPSE! 🚨")
			}
			h.sessionMgr.Delete(userID)
			return h.sendWithKeyboard(chatID,
				fmt.Sprintf("💥 TARGET DESTROYED! DEBT \"%s\" COMPLETELY ERADICATED!", state.TempDebt.Description),
				h.debtsKeyboard())
		}

		state.TempDebt.Amount = newAmount
		if err := h.storage.Update(h.ctx, state.TempDebt); err != nil {
			h.sessionMgr.Delete(userID)
			h.logger.Errorf("failed to update debt %d: %v", state.TempDebt.ID, err)
			return h.sendErrorMessage(chatID, "DEBT WEAKENING FAILED! SPIRAL ENERGY INSUFFICIENT! 🚨")
		}

		h.sessionMgr.Delete(userID)
		return h.sendWithKeyboard(chatID,
			fmt.Sprintf("🚀 PAYMENT TORPEDO LAUNCHED! TARGET DAMAGED!\n\n"+
				"▫️ COSMIC DEBT ENTITY: %s\n"+
				"▫️ REMAINING MASS: %s₽\n\n"+
				"💢 INITIATE SECONDARY ATTACK RUN?",
				state.TempDebt.Description,
				formatMoney(newAmount)),
			h.debtsKeyboard())

	default:
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("UNKNOWN DRILL SEQUENCE: %s 🚨", step))
	}
}
