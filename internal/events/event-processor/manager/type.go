package manager

import (
	"encoding/json"
	"fmt"
	"time"

	"drillCore/internal/model"
	"drillCore/internal/session"
)

type TypeHandler int

const (
	IgnoreHandler TypeHandler = iota
	CMDHandler
	DateHandler
	MainMenuHandler
	DebtHandler
)

type Step int

const (
	StepIgnore Step = iota
	StepStart
	StepList
	StepAddStart
	StepSelect
	StepEditStart
	StepPayStart
	StepEnterPayment
	StepPayAmount
	StepPayFinish
	StepYear
	StepMonth
	StepDay
	StepDeleteStart
	StepAddAmount
	StepAddDescription
	StepEditMenu
	StepEnterDate
	StepEditDate
	StepEnterAmount
	StepEditAmount
	StepEnterDescription
	StepEditDescription
	StepEditFinish
	StepDeleteConfirm
	StepDeleteFinish
	StepAddFinish
)

type State struct {
	BackHandler TypeHandler
	BackStep    Step

	Handler TypeHandler
	Step    Step

	NextHandler TypeHandler
	NextStep    Step

	TempDebt *model.Debt
	TempDate *time.Time
}

func ExtractState(session *session.Session) (*State, error) {
	if state, ok := session.State.(State); ok {
		return &state, nil
	}

	if statePtr, ok := session.State.(*State); ok {
		return statePtr, nil
	}

	return nil, fmt.Errorf("failed to extract state: expected State or *State, got %T", session.State)
}

type CallBack struct {
	Handler TypeHandler `json:"h"`
	Step    Step        `json:"s"`
	Data    string      `json:"d"`
}

// CreateCallBack — serializes callback data into compact JSON array
func CreateCallBack(handler TypeHandler, step Step, data string) (string, error) {
	cb := CallBack{
		Handler: handler,
		Step:    step,
		Data:    data,
	}
	jsonData, err := json.Marshal(cb)
	if err != nil {
		return "", fmt.Errorf("failed to marshal callback: %w", err)
	}

	return string(jsonData), nil
}

// ParseCallBack — deserializes callback data
func ParseCallBack(data string) (*CallBack, error) {
	var cb CallBack
	if err := json.Unmarshal([]byte(data), &cb); err != nil {
		return nil, fmt.Errorf("invalid callback data: %w, raw: %s", err, data)
	}
	return &cb, nil
}
