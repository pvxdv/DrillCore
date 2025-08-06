package debt

import (
	"fmt"
	"strconv"
	"time"

	"drillCore/internal/bot"
	"drillCore/internal/events/event-processor/manager"
	"drillCore/internal/model"
)

func (h *Handler) menuKeyboard() (bot.ReplyMarkup, error) {
	addCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepAddStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	editCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepEditStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	payCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepPayStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	deleteCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepDeleteStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	listCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepList, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	mainMenuCb, err := manager.CreateCallBack(manager.MainMenuHandler, manager.StepStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}
	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{
			{Text: manager.AddDebtButton, CallbackData: addCb},
		},
		{
			{Text: manager.EditDebtButton, CallbackData: editCb},
		},
		{
			{Text: manager.PayDebtButton, CallbackData: payCb},
		},
		{
			{Text: manager.DeleteDebtButton, CallbackData: deleteCb},
		},
		{
			{Text: manager.ListDebtButton, CallbackData: listCb},
		},
		{
			{Text: manager.MainMenuButton, CallbackData: mainMenuCb},
		},
	}), nil
}

func (h *Handler) editKeyboard() (bot.ReplyMarkup, error) {
	descriptionEditCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepEnterDescription, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	amountEditCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepEnterAmount, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	dateEditCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepEnterDate, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	confirmEditCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepEditFinish, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	cancelCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{
			{Text: manager.EditDescButton, CallbackData: descriptionEditCb},
		},
		{
			{Text: manager.EditAmountButton, CallbackData: amountEditCb},
		},
		{
			{Text: manager.EditDateButton, CallbackData: dateEditCb},
		},
		{
			{Text: manager.ConfirmEditButton, CallbackData: confirmEditCb},
		},
		{
			{Text: manager.CancelButton, CallbackData: cancelCb},
		},
	}), nil
}

func (h *Handler) selectKeyboard(debts []*model.Debt) (bot.ReplyMarkup, error) {
	cancelCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	var buttons [][]bot.InlineKeyboardButton
	for _, d := range debts {
		btnText := fmt.Sprintf("🌀 %s - %s₽",
			truncate(d.Description, 20),
			formatMoney(d.Amount))

		if d.ReturnDate != nil && d.ReturnDate.Before(time.Now()) {
			btnText = fmt.Sprintf("💢 %s - %s₽",
				truncate(d.Description, 20),
				formatMoney(d.Amount))
		}

		selectCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepSelect, strconv.FormatInt(d.ID, 10))
		if err != nil {
			return bot.ReplyMarkup{}, err
		}

		buttons = append(buttons, []bot.InlineKeyboardButton{
			{
				Text:         btnText,
				CallbackData: selectCb,
			},
		})
	}

	buttons = append(buttons, []bot.InlineKeyboardButton{
		{Text: manager.CancelButton, CallbackData: cancelCb},
	})

	return bot.NewInlineKeyboard(buttons), nil
}

func (h *Handler) confirmKeyboard(confirmStep manager.Step) (bot.ReplyMarkup, error) {
	confirmCb, err := manager.CreateCallBack(manager.DebtHandler, confirmStep, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	cancelCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{
			{Text: manager.ConfirmButton, CallbackData: confirmCb},
			{Text: manager.CancelButton, CallbackData: cancelCb},
		},
	}), nil
}

func (h *Handler) dateKeyboard(backH manager.TypeHandler, backS manager.Step) (bot.ReplyMarkup, error) {
	dateCb, err := manager.CreateCallBack(manager.DateHandler, manager.StepYear, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	backCb, err := manager.CreateCallBack(backH, backS, "")
	if err != nil {
		h.logger.Errorf("failed to create back button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	cancelCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{
			{Text: manager.SelectDateButton, CallbackData: dateCb},
		},
		{
			{Text: manager.BackStepButton, CallbackData: backCb},
		},
		{
			{Text: manager.CancelButton, CallbackData: cancelCb},
		},
	}), nil
}

func (h *Handler) cancelKeyboard() (bot.ReplyMarkup, error) {
	cancelCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{{Text: manager.CancelButton, CallbackData: cancelCb}},
	}), nil
}

func (h *Handler) redirectKeyboard(backH manager.TypeHandler, backS manager.Step, nextH manager.TypeHandler, nextS manager.Step) (bot.ReplyMarkup, error) {
	backCb, err := manager.CreateCallBack(backH, backS, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	confirmCb, err := manager.CreateCallBack(nextH, nextS, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	cancelCb, err := manager.CreateCallBack(manager.DebtHandler, manager.StepStart, "")
	if err != nil {
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{
			{Text: manager.RedirectDebtButton, CallbackData: confirmCb},
		},
		{
			{Text: manager.BackStepButton, CallbackData: backCb},
		},
		{
			{Text: manager.CancelButton, CallbackData: cancelCb},
		},
	}), nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
