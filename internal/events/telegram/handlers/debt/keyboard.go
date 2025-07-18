package debt

import (
	tgClient "drillCore/internal/clients/telergam"
	"drillCore/internal/model"
	"fmt"
	"sort"

	"time"
)

func (h *Handler) debtsKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: "🌪️ DRILL THROUGH", CallbackData: cbAddStart},
			{Text: "🌀 RESHAPE", CallbackData: cbEditStart},
		},
		{
			{Text: "💥 SMASH DEBT", CallbackData: cbPayStart},
			{Text: "☠️ ANNIHILATE", CallbackData: cbDeleteStart},
		},
		{
			{Text: "📜 BATTLE LOG", CallbackData: cbList},
		},
		{
			{Text: "🌀 SPIRAL MENU", CallbackData: "main_menu"},
		},
	})
}

func (h *Handler) confirmPayKeyboard(amount int64) tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: "✅ Подтвердить", CallbackData: cbPayConfirm + fmt.Sprintf("_%d", amount)},
			{Text: "❌ Отменить", CallbackData: cbCancel},
		},
	})
}

func (h *Handler) cancelKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{{Text: "❌ Отмена", CallbackData: cbCancel}},
	})
}

func (h *Handler) debtsListKeyboard(debts []*model.Debt, flow string) tgClient.ReplyMarkup {
	// Сортируем: сначала активные (по дате), затем просроченные
	sort.Slice(debts, func(i, j int) bool {
		now := time.Now()

		// Оба долга с датой
		if debts[i].ReturnDate != nil && debts[j].ReturnDate != nil {
			iOverdue := debts[i].ReturnDate.Before(now)
			jOverdue := debts[j].ReturnDate.Before(now)

			// Если оба просрочены - сортируем по степени просрочки
			if iOverdue && jOverdue {
				return debts[i].ReturnDate.Before(*debts[j].ReturnDate)
			}
			// Просроченные всегда ниже
			if iOverdue {
				return false
			}
			if jOverdue {
				return true
			}
			// Оба не просрочены - сортируем по дате
			return debts[i].ReturnDate.Before(*debts[j].ReturnDate)
		}

		// Долги без даты идут после долгов с датой
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
		// Форматируем текст кнопки
		btnText := fmt.Sprintf("▫️ %s - %s₽",
			truncate(d.Description, 20),
			formatMoney(d.Amount))

		// Помечаем просроченные
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

	// Добавляем кнопку "Назад"
	buttons = append(buttons, []tgClient.InlineKeyboardButton{
		{Text: "↩️ Назад", CallbackData: cbMenu},
	})

	return tgClient.NewInlineKeyboard(buttons)
}

// Сокращение длинного текста
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func (h *Handler) dateKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: "Выбрать дату", CallbackData: cbDateStart},
		},
		{
			{Text: "❌ Отмена", CallbackData: cbCancel},
		},
	})
}

func (h *Handler) yearKeyboard() tgClient.ReplyMarkup {
	currentYear := time.Now().Year()
	var rows [][]tgClient.InlineKeyboardButton

	// Первый ряд: текущий и следующие 2 года
	var firstRow []tgClient.InlineKeyboardButton
	for y := currentYear; y <= currentYear+2; y++ {
		firstRow = append(firstRow, tgClient.InlineKeyboardButton{
			Text:         fmt.Sprintf("📅 %d", y),
			CallbackData: fmt.Sprintf(cbYear, y),
		})
	}
	rows = append(rows, firstRow)

	// Второй ряд: следующие 3 года
	var secondRow []tgClient.InlineKeyboardButton
	for y := currentYear + 3; y <= currentYear+5; y++ {
		secondRow = append(secondRow, tgClient.InlineKeyboardButton{
			Text:         fmt.Sprintf("📅 %d", y),
			CallbackData: fmt.Sprintf(cbYear, y),
		})
	}
	rows = append(rows, secondRow)

	// Кнопка отмены
	rows = append(rows, []tgClient.InlineKeyboardButton{
		{Text: "❌ Отмена", CallbackData: cbCancel},
	})

	return tgClient.NewInlineKeyboard(rows)
}

func (h *Handler) monthKeyboard() tgClient.ReplyMarkup {
	var rows [][]tgClient.InlineKeyboardButton

	for m := 1; m <= 12; {
		var row []tgClient.InlineKeyboardButton
		for i := 0; i < 4 && m <= 12; i++ { // Исправлено условие
			row = append(row, tgClient.InlineKeyboardButton{
				Text:         time.Month(m).String(), // Лучше показывать названия месяцев
				CallbackData: fmt.Sprintf(cbMonth, m),
			})
			m++
		}
		rows = append(rows, row)
	}

	rows = append(rows, []tgClient.InlineKeyboardButton{
		{Text: "❌ Отмена", CallbackData: cbCancel},
	})

	return tgClient.NewInlineKeyboard(rows)
}

func (h *Handler) dayKeyboard(month int) tgClient.ReplyMarkup {
	now := time.Now()
	year := now.Year() // Используем текущий год
	daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	var rows [][]tgClient.InlineKeyboardButton

	// Заголовок с названием месяца и года
	monthName := time.Month(month).String()
	rows = append(rows, []tgClient.InlineKeyboardButton{
		{Text: fmt.Sprintf("%s %d", monthName, year), CallbackData: "ignore"},
	})

	// Дни недели (заголовки)
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

	// Вычисляем день недели для первого дня месяца
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7 // Воскресенье -> 7
	}

	// Вычисляем номер недели для первого дня
	_, firstWeek := firstDay.ISOWeek()

	currentWeek := firstWeek
	dayCounter := 1

	for dayCounter <= daysInMonth {
		var row []tgClient.InlineKeyboardButton

		// Добавляем номер недели
		row = append(row, tgClient.InlineKeyboardButton{
			Text:         fmt.Sprintf("(%d)", currentWeek),
			CallbackData: "ignore",
		})

		// Добавляем пустые кнопки для выравнивания (для первой недели)
		if dayCounter == 1 {
			for i := 1; i < weekday; i++ {
				row = append(row, tgClient.InlineKeyboardButton{
					Text:         " ",
					CallbackData: "ignore",
				})
			}
		}

		// Добавляем дни месяца
		for len(row) < 8 && dayCounter <= daysInMonth { // 1 (неделя) + 7 дней
			row = append(row, tgClient.InlineKeyboardButton{
				Text:         fmt.Sprintf("%d", dayCounter),
				CallbackData: fmt.Sprintf(cbDay, dayCounter),
			})
			dayCounter++

			// Проверяем смену недели
			if len(row) == 8 { // Если строка заполнена
				currentDate := time.Date(year, time.Month(month), dayCounter, 0, 0, 0, 0, time.UTC)
				_, currentWeek = currentDate.ISOWeek()
			}
		}

		// Заполняем оставшиеся ячейки пустыми кнопками (для последней недели)
		for len(row) < 8 {
			row = append(row, tgClient.InlineKeyboardButton{
				Text:         " ",
				CallbackData: "ignore",
			})
		}

		rows = append(rows, row)
	}

	// Кнопки управления
	rows = append(rows, []tgClient.InlineKeyboardButton{
		{Text: "❌ Отмена", CallbackData: cbCancel},
	})

	return tgClient.NewInlineKeyboard(rows)
}

func (h *Handler) editOptionsKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: "Описание", CallbackData: cbEditDescProcess},
			{Text: "Сумму", CallbackData: cbEditAmountProcess},
			{Text: "Дату", CallbackData: cbEditDateProcess},
		},
		{
			{Text: "✅ Завершить", CallbackData: cbFinishEdit},
			{Text: "❌ Отмена", CallbackData: cbCancel},
		},
	})
}

func (h *Handler) confirmDeleteKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: "✅ Да, удалить", CallbackData: cbDeleteConfirm},
			{Text: "❌ Нет, отменить", CallbackData: cbCancel},
		},
	})
}
