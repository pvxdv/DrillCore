package command

import tgClient "drillCore/internal/clients/telergam"

func (h *Handler) mainKeyboard() tgClient.ReplyMarkup {
	return tgClient.NewInlineKeyboard([][]tgClient.InlineKeyboardButton{
		{{Text: mainButton, CallbackData: cbMainMenu}},
	})
}
