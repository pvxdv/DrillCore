package bot

func NewInlineKeyboard(buttons [][]InlineKeyboardButton) ReplyMarkup {
	return ReplyMarkup{
		InlineKeyboard: buttons,
	}
}
