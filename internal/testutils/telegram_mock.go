package testutils

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MockTelegramBot implements telegram.BotInterface
type MockTelegramBot struct {
	sentMessages     []string
	expectedMessages []string
}

func NewMockTelegramBot() *MockTelegramBot {
	return &MockTelegramBot{
		sentMessages:     make([]string, 0),
		expectedMessages: make([]string, 0),
	}
}

func (m *MockTelegramBot) GetUpdates() ([]tgbotapi.Update, error) {
	return []tgbotapi.Update{}, nil
}

func (m *MockTelegramBot) SendMessage(chatID int64, text string) error {
	m.sentMessages = append(m.sentMessages, text)
	return nil
}

func (m *MockTelegramBot) VerifyMessageSent(t *testing.T, expectedText string) {
	for _, msg := range m.sentMessages {
		if msg == expectedText {
			return
		}
	}
	t.Errorf("Expected message '%s' to be sent", expectedText)
}

func (m *MockTelegramBot) Reset() {
	m.sentMessages = make([]string, 0)
	m.expectedMessages = make([]string, 0)
}

func (m *MockTelegramBot) ExpectMessage(text string) {
	m.expectedMessages = append(m.expectedMessages, text)
}

func (m *MockTelegramBot) VerifyExpectations(t *testing.T) {
	if len(m.expectedMessages) != len(m.sentMessages) {
		t.Errorf("Expected %d messages, got %d", len(m.expectedMessages), len(m.sentMessages))
	}
	// ... verify each message ...
}
