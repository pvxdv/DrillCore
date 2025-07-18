package debt

import (
	"drillCore/internal/model"
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
			h.logger.Errorf("failed to get debts: %v", err)
			return h.sendErrorMessage(chatID, "COSMIC DEBT RADAR MALFUNCTION! 🚨")
		}

		if len(debts) == 0 {
			return h.sendWithKeyboard(chatID, "🌟 YOUR DEBT FIELD IS PRISTINE! NO TARGETS FOR EDIT DRILL!", h.debtsKeyboard())
		}

		// Отправляем список долгов с кнопками для выбора
		return h.sendWithKeyboard(chatID,
			"🌀 SELECT TARGET FOR REALITY EDITING DRILL:",
			h.debtsListKeyboard(debts, flowEdit))
	}

	state, ok := session.State.(*State)
	if !ok {
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, "SPIRAL MEMORY CORRUPTION! SEQUENCE ABORTED! 🚨")
	}

	h.logger.Debugf("get temp debt:%+v", state.TempDebt)

	switch step {
	case "select":
		h.logger.Debugf("handle edit flow step: select, data:%s", data)
		debtID, err := strconv.ParseInt(strings.TrimPrefix(data, "debt_edit_select_"), 10, 64)
		if err != nil {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "INVALID TARGET IDENTIFIER! 🚨")
		}

		debt, err := h.storage.Get(h.ctx, debtID)
		if err != nil {
			h.sessionMgr.Delete(userID)
			h.logger.Errorf("failed to get debt %d: %v", debtID, err)
			return h.sendErrorMessage(chatID, "TARGET LOCK FAILED! SPIRAL SIGNAL LOST! 🚨")
		}

		if debt.UserID != int64(userID) {
			h.sessionMgr.Delete(userID)
			return h.sendErrorMessage(chatID, "💢 DRILL COLLISION DETECTED! THIS DEBT CORE BELONGS TO ANOTHER PILOT! "+
				"YOUR DRILL CANNOT PIERCE ANOTHER MAN'S SOUL! ⚔️")
		}

		state.TempDebt = debt
		state.Step = ""
		h.sessionMgr.Set(userID, h.ID(), state)

		return h.sendWithKeyboard(chatID, "🌀 DRILL EDITING MODE ENGAGED!\n"+
			"WHAT PART OF THIS DEBT CORE SHALL WE DRILL INTO?", h.editOptionsKeyboard())
	case "choice":
		h.logger.Debugf("handle edit flow step: choice, data:%s", data)
		return h.sendWithKeyboard(chatID, "⚡ DRILL BIT SELECTION REQUIRED!\n"+
			"CHOOSE TARGET COMPONENT FOR DEEP DRILLING:", h.editOptionsKeyboard())

	case stepDescription:
		h.logger.Debugf("handle edit flow step: description, data:%s", data)

		if data == "debt_edit_description_process" {
			state.Step = stepDescription
			h.sessionMgr.Set(userID, h.ID(), state)

			return h.sendWithKeyboard(chatID,
				"💢 INITIATE DESCRIPTION DRILLING SEQUENCE!\n"+
					"INPUT NEW TARGET DESIGNATION FOR DEEP CORE:",
				h.cancelKeyboard())
		}

		if strings.TrimSpace(data) == "" {
			return h.sendWithKeyboard(chatID,
				"☄️ DRILL BIT JAMMED! EMPTY DESIGNATION DETECTED!\n"+
					"A DRILL THAT DOESN'T PIERCE IS NO DRILL AT ALL! TRY HARDER:",
				h.cancelKeyboard())
		}

		state.TempDebt.Description = data
		state.Step = ""
		h.sessionMgr.Set(userID, h.ID(), state)

		h.logger.Debugf("update temp debt with description: %+v:", state.TempDebt)

		return h.sendWithKeyboard(chatID,
			"🌌 CORE DRILLING SUCCESS!\n"+
				"NEW TARGET MARKINGS: "+strings.ToUpper(data)+"\n\n"+
				"CONTINUE DRILLING OPERATIONS?",
			h.editOptionsKeyboard())

	case stepAmount:
		h.logger.Debugf("handle edit flow step: amount, data:%s", data)

		if data == "debt_edit_amount_process" {
			state.Step = stepAmount
			h.sessionMgr.Set(userID, h.ID(), state)

			return h.sendWithKeyboard(chatID,
				"💥 INITIATE QUANTUM DRILLING!\n"+
					"INPUT NEW SPIRAL POWER OUTPUT (MIN 1 DRILL UNIT):",
				h.cancelKeyboard())
		}

		if strings.TrimSpace(data) == "" {
			return h.sendWithKeyboard(chatID,
				"💢 DRILL BIT JAMMED! ENERGY INPUT EMPTY! ⚡\n"+
					"MY DRILL IS THE DRILL THAT CREATES NUMBERS!\n"+
					"INPUT SPIRAL POWER VALUE:",
				h.cancelKeyboard())
		}

		amount, err := strconv.ParseInt(data, 10, 64)
		if err != nil || amount <= 0 {
			return h.sendWithKeyboard(chatID,
				"🚨 DRILL OVERLOAD! ONLY POSITIVE NUMBERS CAN PIERCE THE HEAVENS!\n"+
					"REENTER DRILL POWER OUTPUT:",
				h.cancelKeyboard())
		}

		state.TempDebt.Amount = amount
		state.Step = ""
		h.sessionMgr.Set(userID, h.ID(), state)

		h.logger.Debugf("update temp debt with amount: %+v:", state.TempDebt)

		return h.sendWithKeyboard(chatID,
			"⚡ SPIRAL POWER RECALIBRATED TO "+
				formatMoney(amount)+
				" DRILL UNITS!\n\n"+
				"SELECT NEXT DRILLING TARGET:",
			h.editOptionsKeyboard())

	case stepDate:
		h.logger.Debugf("handle edit flow step: date, data:%s", data)
		if data == "debt_edit_date_process" {
			state.Step = stepYear
			h.sessionMgr.Set(userID, h.ID(), state)

			return h.sendWithKeyboard(chatID,
				"⏳ TEMPORAL DRILLING SEQUENCE ENGAGED!\n"+
					"SET DRILL PIERCING COORDINATES:",
				h.dateKeyboard())
		}

		if data == "" {
			state.Step = stepDate
			h.sessionMgr.Set(userID, h.ID(), state)

			return h.sendWithKeyboard(
				chatID,
				"⌛ TIME DRILL SPINNING!\n"+
					"ADJUST TEMPORAL PENETRATION DEPTH:",
				h.dateKeyboard(),
			)
		}
	case "finish":
		h.logger.Debugf("handle edit flow step: finish, data:%s", data)

		h.logger.Debugf("try to save:%+v", state.TempDebt)

		err := h.storage.Update(h.ctx, state.TempDebt)
		h.sessionMgr.Delete(userID)

		if err != nil {
			h.logger.Errorf("failed to update debt: %v", err)

			return h.sendErrorMessage(
				chatID,
				"CATASTROPHIC DRILL FAILURE! THE UNIVERSE RESISTED OUR PIERCING! 🚨",
			)
		}

		return h.sendWithKeyboard(
			chatID,
			"✨ ULTRA DRILLING SEQUENCE COMPLETE!\n\n"+
				"▫️ TARGET: "+
				strings.ToUpper(state.TempDebt.Description)+
				"\n"+
				"▫️ SPIRAL OUTPUT: "+
				formatMoney(state.TempDebt.Amount)+
				" DRILL UNITS\n\n"+
				"💢 WHO THE HELL DO YOU THINK WE ARE?! OUR DRILL PIERCED THROUGH!",
			h.debtsKeyboard(),
		)

	default:
		return h.sendErrorMessage(
			chatID,
			"🚨 UNKNOWN DRILLING SEQUENCE DETECTED!\n"+
				"ABNORMAL DRILL PATTERN IN STEP: "+
				step,
		)
	}

	return h.sendErrorMessage(
		chatID,
		"🚨 UNKNOWN DRILLING SEQUENCE DETECTED!\n"+
			"ABNORMAL DRILL PATTERN IN STEP: "+
			step,
	)
}
