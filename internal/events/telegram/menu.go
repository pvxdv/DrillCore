package telegram

func (p *Processor) showMainMenu(chatID int) error {
	delete(p.userSessions, chatID)
	return p.tg.SendMessageWithKeyboard(chatID, msgWelcome, mainKeyboard)
}
