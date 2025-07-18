package tgClient

func NewInlineKeyboard(buttons [][]InlineKeyboardButton) ReplyMarkup {
	return ReplyMarkup{
		InlineKeyboard: buttons,
	}
}
