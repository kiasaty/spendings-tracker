package app

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasaty/spendings-tracker/internal/testutils"
)

func TestHandleUpdate(t *testing.T) {
	tests := []struct {
		name          string
		update        *tgbotapi.Update
		expectedCost  float64
		expectedTags  []string
		expectedError bool
	}{
		{
			name:          "valid expense with tags",
			update:        testutils.NewTestUpdate(1, 123456789, "Lunch 15.50 #food #work"),
			expectedCost:  15.50,
			expectedTags:  []string{"food", "work"},
			expectedError: false,
		},
		{
			name:          "valid expense without tags",
			update:        testutils.NewTestUpdate(2, 123456789, "Coffee 3.50"),
			expectedCost:  3.50,
			expectedTags:  []string{},
			expectedError: false,
		},
		{
			name:          "expense without price",
			update:        testutils.NewTestUpdate(3, 123456789, "Just a message without price"),
			expectedError: true,
		},
		{
			name:          "expense with invalid price",
			update:        testutils.NewTestUpdate(4, 123456789, "Invalid price abc #food"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks with initial state
			db := testutils.NewMockDatabaseClient()
			bot := testutils.NewMockTelegramBot()

			// Create app with mocks
			app := &App{
				DB:  db,
				Bot: bot,
			}

			// Process the update
			app.handleUpdate(tt.update)

			if tt.expectedError {
				// For error cases, verify no spending was created
				if spending, _ := db.FindSpendingByMessageId(tt.update.Message.MessageID); spending != nil {
					t.Errorf("Expected no spending to be created for error case")
				}
				return
			}

			// Verify spending was created with correct values
			db.VerifySpending(t, tt.update.Message.MessageID, tt.expectedCost)

			// Verify tags were created
			for _, tagName := range tt.expectedTags {
				tag, err := db.FindTagByName(tagName)
				if err != nil {
					t.Errorf("Unexpected error finding tag: %v", err)
				}
				if tag == nil {
					t.Errorf("Expected tag %s to be created", tagName)
				}
			}

			// Reset mocks for next test
			db.Reset()
			bot.Reset()
		})
	}
}

func TestHandleUpdateWithError(t *testing.T) {
	// Test error handling when database operations fail
	db := testutils.NewMockDatabaseClient()
	bot := testutils.NewMockTelegramBot()
	app := &App{
		DB:  db,
		Bot: bot,
	}

	// Configure mock to return error on create
	db.SetErrorOnCreate(true)

	// Try to create a spending
	update := testutils.NewTestUpdate(1, 123456789, "Lunch 15.50")
	app.handleUpdate(update)

	// Verify no spending was created due to error
	if spending, _ := db.FindSpendingByMessageId(1); spending != nil {
		t.Errorf("Expected no spending to be created when database returns error")
	}
}
