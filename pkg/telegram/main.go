package telegram

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Bot *tgbotapi.BotAPI
}

func NewTelegramBot() (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = true

	return &TelegramBot{
		Bot: bot,
	}, nil
}

func (t *TelegramBot) GetUpdates() ([]tgbotapi.Update, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.Bot.GetUpdatesChan(u)
	var result []tgbotapi.Update

	// Get all available updates
	for update := range updates {
		result = append(result, update)
		if len(result) >= 100 { // Limit the number of updates to process at once
			break
		}
	}

	return result, nil
}

func (t *TelegramBot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
