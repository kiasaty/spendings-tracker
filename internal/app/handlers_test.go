package app

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasaty/spendings-tracker/internal/testutils"
)

func TestHandleUpdate(t *testing.T) {
	tests := []struct {
		name         string
		update       *tgbotapi.Update
		expectedCost float64
		expectedTags []string
		expectedDate time.Time
		expectError  bool
	}{
		{
			name:         "Valid expense with tags",
			update:       testutils.NewTestUpdate(1, 123456789, "Lunch 15.50 #food #work"),
			expectedCost: 15.50,
			expectedTags: []string{"food", "work"},
			expectedDate: time.Now(),
		},
		{
			name:         "Valid expense without tags",
			update:       testutils.NewTestUpdate(2, 123456789, "Dinner 25.75"),
			expectedCost: 25.75,
			expectedTags: []string{},
			expectedDate: time.Now(),
		},
		{
			name:         "Valid expense with date",
			update:       testutils.NewTestUpdate(3, 123456789, "Lunch 15.50 2024-05-09 #food"),
			expectedCost: 15.50,
			expectedTags: []string{"food"},
			expectedDate: time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "Valid expense with different date format",
			update:       testutils.NewTestUpdate(4, 123456789, "Dinner 25.75 09.05.2024 #food"),
			expectedCost: 25.75,
			expectedTags: []string{"food"},
			expectedDate: time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "Expense without price",
			update:      testutils.NewTestUpdate(5, 123456789, "Lunch #food"),
			expectError: true,
		},
		{
			name:        "Expense with invalid price",
			update:      testutils.NewTestUpdate(6, 123456789, "Lunch abc #food"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := testutils.NewMockDatabaseClient()
			mockBot := testutils.NewMockTelegramBot()
			app, err := NewApp(mockDB, mockBot)
			if err != nil {
				t.Fatalf("Failed to create app: %v", err)
			}

			app.handleUpdate(tt.update)

			if tt.expectError {
				// Verify no spending was created
				spending, _ := mockDB.FindSpendingByMessageId(tt.update.Message.MessageID)
				if spending != nil {
					t.Errorf("Expected no spending to be created for invalid message")
				}
				return
			}

			// Get the spending once
			spending, _ := mockDB.FindSpendingByMessageId(tt.update.Message.MessageID)

			// Verify spending cost and date
			mockDB.VerifySpending(t, spending, tt.expectedCost, tt.expectedDate)

			// Verify tags
			mockDB.VerifySpendingTags(t, spending, tt.expectedTags)

			// Reset mocks for next test
			mockDB.Reset()
			mockBot.Reset()
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
