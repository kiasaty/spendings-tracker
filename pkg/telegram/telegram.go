package telegram

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotInterface defines the interface for interacting with Telegram
type BotInterface interface {
	GetUpdates() tgbotapi.UpdatesChannel
	SendMessage(chatID int64, text string) error
}

// telegramBot implements the TelegramBot interface
type telegramBot struct {
	bot *tgbotapi.BotAPI
}

// NewTelegramBot creates a new instance of TelegramBot
func NewTelegramBot() (BotInterface, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = true

	return &telegramBot{
		bot: bot,
	}, nil
}

// GetUpdates retrieves updates from Telegram
func (t *telegramBot) GetUpdates() tgbotapi.UpdatesChannel {
	fmt.Println("Setting up update channel...")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return t.bot.GetUpdatesChan(u)
}

// SendMessage sends a message to a Telegram chat
func (t *telegramBot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
