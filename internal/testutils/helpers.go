package testutils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewTestUpdate creates a test update with specific message
func NewTestUpdate(messageID int, chatID int64, text string) *tgbotapi.Update {
	return &tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: messageID,
			Chat:      &tgbotapi.Chat{ID: chatID},
			Text:      text,
		},
	}
}
