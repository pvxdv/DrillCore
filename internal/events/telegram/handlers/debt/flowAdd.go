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

		return h.sendWithKeyboard(chatID,
			"🌀 INITIATE DEBT DRILL DEPLOYMENT SEQUENCE!\n"+
				"ENTER DESCRIPTION:\n",
			h.cancelKeyboard())
	}

	state, ok := session.State.(*State)
	if !ok {
		h.sessionMgr.Delete(userID)
		return h.sendErrorMessage(chatID, fmt.Sprintf("DRILL CORE MALFUNCTION! UNKNOWN STEP: %s 🚨", step))
	}

	switch state.Step {
	case stepDescription:
		if strings.TrimSpace(data) == "" {
			return h.sendWithKeyboard(chatID,
				"💢 DRILL BIT ERROR! EMPTY DESIGNATION DETECTED!\n"+
					"A DRILL NEEDS A TARGET TO PIERCE!\n"+
					"REENTER TARGET DESIGNATION:",
				h.cancelKeyboard())
		}

		state.TempDebt.Description = data
		state.Step = stepAmount
		h.sessionMgr.Set(userID, h.ID(), state)

		return h.sendWithKeyboard(chatID,
			"⚡ TARGET ACQUIRED! NOW INPUT SPIRAL DEBT LOAD:\n"+
				"INPUT AMOUNT",
			h.cancelKeyboard())

	case stepAmount:
		amount, err := strconv.ParseInt(data, 10, 64)
		if err != nil || amount <= 0 {

			return h.sendWithKeyboard(chatID,
				"🚨 ILLEGAL ENERGY INPUT! ONLY POSITIVE NUMBERS CAN POWER THE DRILL!\n"+
					"REENTER SPIRAL DEBT LOAD:",
				h.cancelKeyboard())
		}

		state.TempDebt.Amount = amount
		state.Step = stepYear

		now := time.Now()
		state.TempDebt.ReturnDate = &now

		h.sessionMgr.Set(userID, h.ID(), state)

		return h.sendWithKeyboard(chatID,
			"⏳ TEMPORAL DRILLING SEQUENCE ENGAGED!\n"+
				"SET D-DAY FOR DEBT RECLAMATION:",
			h.dateKeyboard())

	default:
		return h.sendErrorMessage(chatID,
			"💥 UNKNOWN DRILLING PHASE DETECTED!\n"+
				"ABNORMAL SEQUENCE AT STEP: "+step)
	}
}
