package mainmenu

import (
	tgClient "drillCore/internal/clients/telergam"
)

func (h *Handler) mainKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{
			{Text: debtButton, CallbackData: cbDebt},
		},
		// [SPIRAL ENGINE BOOTING...]
		//{
		//    {Text: "🌪️ WEATHER DRILL", CallbackData: cbWeather},
		//    {Text: "⚡ TASK DRILL", CallbackData: cbTask},
		//},
		{
			{Text: mainConsoleButton, CallbackData: "main_menu"},
		},
	})
}
