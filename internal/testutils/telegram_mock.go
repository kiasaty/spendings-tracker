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

func (m *MockTelegramBot) GetUpdates() tgbotapi.UpdatesChannel {
	// For tests, we don't need a real updates channel
	return make(chan tgbotapi.Update)
}

func (m *MockTelegramBot) SendMessage(chatID int64, text string) error {
	m.sentMessages = append(m.sentMessages, text)
	return nil
}

func (m *MockTelegramBot) VerifyMessage(t *testing.T, expectedText string) {
	for _, msg := range m.sentMessages {
		if msg == expectedText {
			return
		}
	}
	t.Errorf("Expected message '%s' to be sent, but got messages: %v", expectedText, m.sentMessages)
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
		return
	}
	for i, expected := range m.expectedMessages {
		if i >= len(m.sentMessages) {
			t.Errorf("Missing expected message: %s", expected)
			continue
		}
		if m.sentMessages[i] != expected {
			t.Errorf("Message mismatch at position %d:\nExpected: %s\nGot: %s", i, expected, m.sentMessages[i])
		}
	}
}
