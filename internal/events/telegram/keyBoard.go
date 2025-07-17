package telegram

import tgClient "drillCore/internal/clients/telergam"

var (
	mainKeyboard = tgClient.ReplyMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{
				{Text: "📋 My Debts", CallbackData: "list_debts"},
				{Text: "➕ Add Debt", CallbackData: "add_debt"},
			},
		},
	}

	backToMainKeyboard = tgClient.ReplyMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{{Text: "🔙 Back to Main", CallbackData: "main_menu"}},
		},
	}
)
