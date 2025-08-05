package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"drillCore/internal/bot"
	"drillCore/internal/events/event-processor/manager"
)

func (h *Handler) dateKeyboard(backH manager.TypeHandler) (bot.ReplyMarkup, error) {
	year, err := manager.CreateCallBack(manager.DateHandler, manager.StepYear, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	back, err := manager.CreateCallBack(backH, manager.StepStart, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{
			{Text: manager.SelectDateButton, CallbackData: year},
		},
		{
			{Text: manager.CancelButton, CallbackData: back},
		},
	}), nil
}

func (h *Handler) yearKeyboard(backH manager.TypeHandler, backS manager.Step) (bot.ReplyMarkup, error) {
	currentYear := time.Now().Year()
	endYear := currentYear + 11
	var rows [][]bot.InlineKeyboardButton

	for year := currentYear; year <= endYear; {
		var row []bot.InlineKeyboardButton

		for i := 0; i < 4 && year <= endYear; i++ {
			callback, err := manager.CreateCallBack(
				manager.DateHandler,
				manager.StepYear,
				strconv.Itoa(year),
			)
			if err != nil {
				h.logger.Errorf("failed to create callback: %v", err)
				return bot.ReplyMarkup{}, err
			}

			row = append(row, bot.InlineKeyboardButton{
				Text:         fmt.Sprintf(manager.SpiralFormat, year),
				CallbackData: callback,
			})
			year++
		}

		rows = append(rows, row)
	}

	backCb, err := manager.CreateCallBack(backH, backS, "")
	if err != nil {
		h.logger.Errorf("failed to create back button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: manager.BackStepButton, CallbackData: backCb},
	})

	cancelCb, err := manager.CreateCallBack(backH, manager.StepStart, "")
	if err != nil {
		h.logger.Errorf("failed to create back button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: manager.CancelButton, CallbackData: cancelCb},
	})

	return bot.NewInlineKeyboard(rows), nil
}

func (h *Handler) monthKeyboard(backH manager.TypeHandler, backS manager.Step) (bot.ReplyMarkup, error) {
	var rows [][]bot.InlineKeyboardButton
	var row []bot.InlineKeyboardButton

	for buttonText, month := range monthButtonMap {
		callbackData, err := manager.CreateCallBack(
			manager.DateHandler,
			manager.StepMonth,
			buttonText,
		)
		if err != nil {
			h.logger.Errorf("failed to create callback for month %s: %v", month, err)
			return bot.ReplyMarkup{}, err
		}

		row = append(row, bot.InlineKeyboardButton{
			Text:         buttonText,
			CallbackData: callbackData,
		})

		if len(row) == 4 {
			rows = append(rows, row)
			row = nil
		}
	}

	if len(row) > 0 {
		rows = append(rows, row)
	}

	backToYear, err := manager.CreateCallBack(manager.DateHandler, manager.StepYear, "")
	if err != nil {
		h.logger.Errorf("failed to create back year button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: manager.ReInputYearButton, CallbackData: backToYear},
	})

	backCb, err := manager.CreateCallBack(backH, backS, "")
	if err != nil {
		h.logger.Errorf("failed to create back step button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: manager.BackStepButton, CallbackData: backCb},
	})

	cancelCb, err := manager.CreateCallBack(backH, manager.StepStart, "")
	if err != nil {
		h.logger.Errorf("failed to create cancel button, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: manager.CancelButton, CallbackData: cancelCb},
	})

	return bot.NewInlineKeyboard(rows), nil
}

func (h *Handler) dayKeyboard(year int, month time.Month, backH manager.TypeHandler, backS manager.Step) (bot.ReplyMarkup, error) {
	var rows [][]bot.InlineKeyboardButton

	// header month + year (inactive)
	monthName := strings.ToUpper(month.String()[:3])

	ignore, err := manager.CreateCallBack(manager.IgnoreHandler, manager.StepIgnore, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	button := fmt.Sprintf("%s %d", monthName, year)

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: button, CallbackData: ignore},
	})

	// header weekdays (inactive)
	var headerRow []bot.InlineKeyboardButton
	headerRow = append(headerRow, bot.InlineKeyboardButton{
		Text:         manager.SpiralEmoji,
		CallbackData: ignore,
	})

	for _, day := range weekdays {
		headerRow = append(headerRow, bot.InlineKeyboardButton{
			Text:         day,
			CallbackData: ignore,
		})
	}

	rows = append(rows, headerRow)

	// days buttons (active)
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	_, firstWeek := firstDay.ISOWeek()

	currentWeek := firstWeek
	dayCounter := 1

	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for dayCounter <= daysInMonth {
		var row []bot.InlineKeyboardButton

		wButton := fmt.Sprintf("%dW", currentWeek)

		row = append(row, bot.InlineKeyboardButton{
			Text:         wButton,
			CallbackData: ignore,
		})

		if dayCounter == 1 {
			for i := 1; i < weekday; i++ {
				row = append(row, bot.InlineKeyboardButton{
					Text:         manager.SpiralEmoji,
					CallbackData: ignore,
				})
			}
		}

		for len(row) < 8 && dayCounter <= daysInMonth {
			dayStr := strconv.Itoa(dayCounter)

			day, err := manager.CreateCallBack(manager.DateHandler, manager.StepDay, dayStr)
			if err != nil {
				h.logger.Errorf("failed to create callback, err:%v", err)
				return bot.ReplyMarkup{}, err
			}

			row = append(row, bot.InlineKeyboardButton{
				Text:         dayStr,
				CallbackData: day,
			})
			dayCounter++

			if len(row) == 8 {
				currentDate := time.Date(year, month, dayCounter, 0, 0, 0, 0, time.UTC)
				_, currentWeek = currentDate.ISOWeek()
			}
		}

		for len(row) < 8 {
			row = append(row, bot.InlineKeyboardButton{
				Text:         manager.SpiralEmoji,
				CallbackData: ignore,
			})
		}

		rows = append(rows, row)
	}

	// navigation buttons
	backToMonth, err := manager.CreateCallBack(manager.DateHandler, manager.StepMonth, "")
	if err != nil {
		h.logger.Errorf("failed to create back to month button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: manager.ReInputMonthButton, CallbackData: backToMonth},
	})

	backCb, err := manager.CreateCallBack(backH, backS, "")
	if err != nil {
		h.logger.Errorf("failed to create back step button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: manager.BackStepButton, CallbackData: backCb},
	})

	cancelCb, err := manager.CreateCallBack(backH, manager.StepStart, "")
	if err != nil {
		h.logger.Errorf("failed to create cancel button calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	rows = append(rows, []bot.InlineKeyboardButton{
		{Text: manager.CancelButton, CallbackData: cancelCb},
	})

	return bot.NewInlineKeyboard(rows), nil
}

func (h *Handler) redirectKeyboard(backH manager.TypeHandler, backS manager.Step, nextH manager.TypeHandler, nextS manager.Step) (bot.ReplyMarkup, error) {
	redirect, err := manager.CreateCallBack(nextH, nextS, "")
	if err != nil {
		h.logger.Errorf("failed to create calldack, err:%v", err)

		return bot.ReplyMarkup{}, err
	}

	backToDay, err := manager.CreateCallBack(manager.DateHandler, manager.StepDay, "")
	if err != nil {
		h.logger.Errorf("failed to create back to day button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	backCb, err := manager.CreateCallBack(backH, backS, "")
	if err != nil {
		h.logger.Errorf("failed to create back step button callback: %v", err)
		return bot.ReplyMarkup{}, err
	}

	cancelCb, err := manager.CreateCallBack(backH, manager.StepStart, "")
	if err != nil {
		h.logger.Errorf("failed to create cancel button calldack, err:%v", err)
		return bot.ReplyMarkup{}, err
	}

	return bot.NewInlineKeyboard([][]bot.InlineKeyboardButton{
		{
			{Text: manager.RedirectDateButton, CallbackData: redirect},
		},
		{
			{Text: manager.ReInputDayButton, CallbackData: backToDay},
		},
		{
			{Text: manager.BackStepButton, CallbackData: backCb},
		},
		{
			{Text: manager.CancelButton, CallbackData: cancelCb},
		},
	}), nil
}
