package telegram

import (
	"context"
	tgClient "drillCore/internal/clients/telergam"
	"drillCore/internal/model"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Debt management methods
func (p *Processor) listDebts(chatID int, userID int) error {
	debts, err := p.storage.Debts(context.TODO(), int64(userID))
	if err != nil {
		log.Printf("Failed to get debts: %v", err)
		return p.sendErrorMessage(chatID, "Failed to retrieve debts")
	}

	if len(debts) == 0 {
		return p.sendWithKeyboard(chatID, msgNoDebts, backToMainKeyboard)
	}

	for _, debt := range debts {
		if err := p.sendDebtDetails(chatID, *debt); err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) sendDebtDetails(chatID int, debt model.Debt) error {
	message := p.formatDebtMessage(debt)
	keyboard := p.createDebtKeyboard(debt.ID)
	return p.sendWithKeyboard(chatID, message, keyboard)
}

func (p *Processor) formatDebtMessage(debt model.Debt) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("🆔 ID: %d\n", debt.ID))
	builder.WriteString(fmt.Sprintf("📝 Description: %s\n", debt.Description))
	builder.WriteString(fmt.Sprintf("💰 Amount: %d\n", debt.Amount))
	if debt.ReturnDate != nil {
		builder.WriteString(fmt.Sprintf("📅 Due Date: %s", debt.ReturnDate.Format("02.01.2006")))
	}
	return builder.String()
}

func (p *Processor) createDebtKeyboard(debtID int64) tgClient.ReplyMarkup {
	return tgClient.ReplyMarkup{
		InlineKeyboard: [][]tgClient.InlineKeyboardButton{
			{
				{Text: "❌ Delete", CallbackData: fmt.Sprintf("delete_debt_%d", debtID)},
			},
			{
				{Text: "🔙 Back", CallbackData: "main_menu"},
			},
		},
	}
}

// Debt creation flow
func (p *Processor) startAddDebt(chatID int) error {
	p.userSessions[chatID] = &SessionState{
		Action:    waitDescription,
		TempDebt:  &model.Debt{},
		MessageID: 0,
	}
	return p.sendWithKeyboard(chatID, msgEnterDesc, backToMainKeyboard)
}

func (p *Processor) handleDebtDescription(chatID int, desc string) error {
	session := p.getSession(chatID)
	session.TempDebt.Description = desc
	session.Action = waitAmount
	return p.sendWithKeyboard(chatID, msgEnterAmount, backToMainKeyboard)
}

func (p *Processor) handleDebtAmount(chatID int, amountStr string) error {
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return p.sendWithKeyboard(chatID, msgInvalidAmount, backToMainKeyboard)
	}

	session := p.getSession(chatID)
	session.TempDebt.Amount = amount
	session.Action = waitDate
	return p.sendWithKeyboard(chatID, msgEnterDate, backToMainKeyboard)
}

func (p *Processor) handleDebtDate(chatID int, dateStr string) error {
	session := p.getSession(chatID)

	if !strings.EqualFold(dateStr, "/skip") {
		date, err := time.Parse("02.01.2006", dateStr)
		if err != nil {
			return p.sendWithKeyboard(chatID, msgInvalidDate, backToMainKeyboard)
		}
		session.TempDebt.ReturnDate = &date
	}

	if _, err := p.storage.Save(context.Background(), session.TempDebt); err != nil {
		log.Printf("Failed to save debt: %v", err)
		return p.sendErrorMessage(chatID, "Failed to save debt")
	}

	delete(p.userSessions, chatID)
	return p.sendWithKeyboard(chatID, msgDebtAdded, mainKeyboard)
}

// Debt deletion
func (p *Processor) deleteDebt(chatID int, debtID int64) error {
	if err := p.storage.Delete(context.Background(), debtID); err != nil {
		log.Printf("Failed to delete debt %d: %v", debtID, err)
		return p.sendErrorMessage(chatID, "Failed to delete debt")
	}
	return p.sendWithKeyboard(chatID, msgDebtDeleted, mainKeyboard)
}

// Helper methods
func (p *Processor) getSession(chatID int) *SessionState {
	return p.userSessions[chatID]
}

func (p *Processor) sendWithKeyboard(chatID int, text string, keyboard tgClient.ReplyMarkup) error {
	if err := p.tg.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		log.Printf("Failed to send message to %d: %v", chatID, err)
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}

func (p *Processor) sendErrorMessage(chatID int, message string) error {
	return p.sendWithKeyboard(chatID, "⚠️ "+message, backToMainKeyboard)
}
