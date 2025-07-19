package debt

import (
	tgClient "drillCore/internal/clients/telergam"
	mainmenu "drillCore/internal/events/telegram/handlers/main-menu"
	"drillCore/internal/model"
	"fmt"
	"sort"

	"time"
)

const (
	addDebtButton    = "💢 NEW DRILL"
	editDebtButton   = "🌀 EDIT DRILL"
	payDebtButton    = "💥 PAY DRILL"
	deleteDebtButton = "☠️ KILL DRILL"
	listDebtButton   = "📜 DRILL LOG"

	confirmButton = "💢 COMMIT DRILL"
	cancelButton  = "🌀 SPIRAL BACK"

	editDescButton   = "🌀 DESC"
	editAmountButton = "💢 CASH"
	editDateButton   = "⏳ TIME"
	finishEditButton = "🌀 FINALIZE DRILLING"

	selectDateButton = "⏳ SET D-DAY"
	yearButtonFormat = "📅 %d"

	ignoreButton = " "
	backButton   = "💢 SPIRAL COMMAND CENTER 💢"
)

func (h *Handler) debtsKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: addDebtButton, CallbackData: cbAddStart},
			{Text: editDebtButton, CallbackData: cbEditStart},
		},
		{
			{Text: payDebtButton, CallbackData: cbPayStart},
			{Text: deleteDebtButton, CallbackData: cbDeleteStart},
		},
		{
			{Text: listDebtButton, CallbackData: cbList},
		},
		{
			{Text: mainmenu.MainConsoleButton, CallbackData: cmdMainMenu},
		},
	})
}

func (h *Handler) confirmPayKeyboard(amount int64) tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: confirmButton, CallbackData: cbPayConfirm + fmt.Sprintf("_%d", amount)},
			{Text: cancelButton, CallbackData: cbCancel},
		},
	})
}

func (h *Handler) cancelKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{{Text: cancelButton, CallbackData: cbCancel}},
	})
}

func (h *Handler) debtsListKeyboard(debts []*model.Debt, flow string) tgClient.ReplyMarkup {
	sort.Slice(debts, func(i, j int) bool {
		now := time.Now()

		if debts[i].ReturnDate != nil && debts[j].ReturnDate != nil {
			iOverdue := debts[i].ReturnDate.Before(now)
			jOverdue := debts[j].ReturnDate.Before(now)

			if iOverdue && jOverdue {
				return debts[i].ReturnDate.Before(*debts[j].ReturnDate)
			}

			if iOverdue {
				return false
			}
			if jOverdue {
				return true
			}

			return debts[i].ReturnDate.Before(*debts[j].ReturnDate)
		}

		if debts[i].ReturnDate == nil {
			return false
		}
		if debts[j].ReturnDate == nil {
			return true
		}
		return true
	})

	var buttons [][]tgClient.InlineKeyboardButton
	for _, d := range debts {
		btnText := fmt.Sprintf("▫️ %s - %s₽",
			truncate(d.Description, 20),
			formatMoney(d.Amount))

		if d.ReturnDate != nil && d.ReturnDate.Before(time.Now()) {
			btnText = fmt.Sprintf("❗ %s - %s₽",
				truncate(d.Description, 20),
				formatMoney(d.Amount))
		}

		buttons = append(buttons, []tgClient.InlineKeyboardButton{
			{
				Text:         btnText,
				CallbackData: fmt.Sprintf(prefix+flow+"_select_%d", d.ID),
			},
		})
	}

	buttons = append(buttons, []tgClient.InlineKeyboardButton{
		{Text: backButton, CallbackData: cbMenu},
	})

	return tgClient.NewInlineKeyboard(buttons)
}

func (h *Handler) dateKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: selectDateButton, CallbackData: cbDateStart},
		},
		{
			{Text: cancelButton, CallbackData: cbCancel},
		},
	})
}

func (h *Handler) yearKeyboard() tgClient.ReplyMarkup {
	currentYear := time.Now().Year()
	var rows [][]tgClient.InlineKeyboardButton

	var firstRow []tgClient.InlineKeyboardButton
	for y := currentYear; y <= currentYear+2; y++ {
		firstRow = append(firstRow, tgClient.InlineKeyboardButton{
			Text:         fmt.Sprintf(yearButtonFormat, y),
			CallbackData: fmt.Sprintf(cbYear, y),
		})
	}
	rows = append(rows, firstRow)

	var secondRow []tgClient.InlineKeyboardButton
	for y := currentYear + 3; y <= currentYear+5; y++ {
		secondRow = append(secondRow, tgClient.InlineKeyboardButton{
			Text:         fmt.Sprintf(yearButtonFormat, y),
			CallbackData: fmt.Sprintf(cbYear, y),
		})
	}
	rows = append(rows, secondRow)

	rows = append(rows, []tgClient.InlineKeyboardButton{
		{Text: cancelButton, CallbackData: cbCancel},
	})

	return tgClient.NewInlineKeyboard(rows)
}

func (h *Handler) monthKeyboard() tgClient.ReplyMarkup {
	var rows [][]tgClient.InlineKeyboardButton

	for m := 1; m <= 12; {
		var row []tgClient.InlineKeyboardButton
		for i := 0; i < 4 && m <= 12; i++ {
			row = append(row, tgClient.InlineKeyboardButton{
				Text:         time.Month(m).String(),
				CallbackData: fmt.Sprintf(cbMonth, m),
			})
			m++
		}
		rows = append(rows, row)
	}

	rows = append(rows, []tgClient.InlineKeyboardButton{
		{Text: cancelButton, CallbackData: cbCancel},
	})

	return tgClient.NewInlineKeyboard(rows)
}

func (h *Handler) dayKeyboard(month int) tgClient.ReplyMarkup {
	now := time.Now()
	year := now.Year()
	daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	var rows [][]tgClient.InlineKeyboardButton

	monthName := time.Month(month).String()
	rows = append(rows, []tgClient.InlineKeyboardButton{
		{Text: fmt.Sprintf("%s %d", monthName, year), CallbackData: "ignore"},
	})

	weekdays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	var headerRow []tgClient.InlineKeyboardButton
	headerRow = append(headerRow, tgClient.InlineKeyboardButton{
		Text:         "№",
		CallbackData: "ignore",
	})
	for _, day := range weekdays {
		headerRow = append(headerRow, tgClient.InlineKeyboardButton{
			Text:         day,
			CallbackData: "ignore",
		})
	}
	rows = append(rows, headerRow)

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	_, firstWeek := firstDay.ISOWeek()

	currentWeek := firstWeek
	dayCounter := 1

	for dayCounter <= daysInMonth {
		var row []tgClient.InlineKeyboardButton

		row = append(row, tgClient.InlineKeyboardButton{
			Text:         fmt.Sprintf("(%d)", currentWeek),
			CallbackData: "ignore",
		})

		if dayCounter == 1 {
			for i := 1; i < weekday; i++ {
				row = append(row, tgClient.InlineKeyboardButton{
					Text:         ignoreButton,
					CallbackData: "ignore",
				})
			}
		}

		for len(row) < 8 && dayCounter <= daysInMonth { // 1 (неделя) + 7 дней
			row = append(row, tgClient.InlineKeyboardButton{
				Text:         fmt.Sprintf("%d", dayCounter),
				CallbackData: fmt.Sprintf(cbDay, dayCounter),
			})
			dayCounter++

			if len(row) == 8 {
				currentDate := time.Date(year, time.Month(month), dayCounter, 0, 0, 0, 0, time.UTC)
				_, currentWeek = currentDate.ISOWeek()
			}
		}

		for len(row) < 8 {
			row = append(row, tgClient.InlineKeyboardButton{
				Text:         " ",
				CallbackData: "ignore",
			})
		}

		rows = append(rows, row)
	}

	rows = append(rows, []tgClient.InlineKeyboardButton{
		{Text: cancelButton, CallbackData: cbCancel},
	})

	return tgClient.NewInlineKeyboard(rows)
}

func (h *Handler) editOptionsKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: editDescButton, CallbackData: cbEditDescProcess},
			{Text: editAmountButton, CallbackData: cbEditAmountProcess},
			{Text: editDateButton, CallbackData: cbEditDateProcess},
		},
		{
			{Text: finishEditButton, CallbackData: cbFinishEdit},
			{Text: cancelButton, CallbackData: cbCancel},
		},
	})
}

func (h *Handler) confirmDeleteKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: confirmButton, CallbackData: cbDeleteConfirm},
			{Text: cancelButton, CallbackData: cbCancel},
		},
	})
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
