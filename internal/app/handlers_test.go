package app

import (
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasaty/spendings-tracker/internal/testutils"
	"github.com/kiasaty/spendings-tracker/models"
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

func TestHandleReportCommand(t *testing.T) {
	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonthStart := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())

	tests := []struct {
		name           string
		command        string
		spendings      []*models.Spending
		expectedReport string
	}{
		{
			name:    "Current month report with multiple tags",
			command: "/report",
			spendings: []*models.Spending{
				{
					MessageId: 1,
					Cost:      15.50,
					SpentAt:   currentMonthStart.AddDate(0, 0, 1),
					Tags:      []models.Tag{{Name: "food"}},
				},
				{
					MessageId: 2,
					Cost:      25.75,
					SpentAt:   currentMonthStart.AddDate(0, 0, 2),
					Tags:      []models.Tag{{Name: "food"}, {Name: "work"}},
				},
				{
					MessageId: 3,
					Cost:      10.00,
					SpentAt:   currentMonthStart.AddDate(0, 0, 3),
					Tags:      []models.Tag{{Name: "work"}},
				},
			},
			expectedReport: "Spending report for current month:\n\nfood: 41.25\nwork: 35.75\n\nTotal: 51.25",
		},
		{
			name:    "Last month report",
			command: "/report_last_month",
			spendings: []*models.Spending{
				{
					MessageId: 4,
					Cost:      30.00,
					SpentAt:   lastMonthStart.AddDate(0, 0, 1),
					Tags:      []models.Tag{{Name: "food"}},
				},
				{
					MessageId: 5,
					Cost:      20.00,
					SpentAt:   lastMonthStart.AddDate(0, 0, 2),
					Tags:      []models.Tag{},
				},
			},
			expectedReport: "Spending report for last month:\n\nfood: 30.00\nother: 20.00\n\nTotal: 50.00",
		},
		{
			name:           "Empty current month report",
			command:        "/report",
			spendings:      []*models.Spending{},
			expectedReport: "Spending report for current month:\n\nTotal: 0.00",
		},
		{
			name:    "Current month report with only untagged spendings",
			command: "/report",
			spendings: []*models.Spending{
				{
					MessageId: 6,
					Cost:      33.00,
					SpentAt:   currentMonthStart.AddDate(0, 0, 1),
					Tags:      []models.Tag{},
				},
			},
			expectedReport: "Spending report for current month:\n\nother: 33.00\n\nTotal: 33.00",
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

			// Add test spendings to mock DB
			for _, spending := range tt.spendings {
				mockDB.CreateSpending(spending)
			}

			// Create update with command
			update := &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: tt.command,
					Chat: &tgbotapi.Chat{
						ID: 123456789,
					},
				},
			}

			// Handle the command
			if tt.command == "/report" {
				app.handleReportCommand(update.Message, false)
			} else {
				app.handleReportCommand(update.Message, true)
			}

			// Verify the report message
			mockBot.VerifyMessage(t, tt.expectedReport)

			// Reset mocks for next test
			mockDB.Reset()
			mockBot.Reset()
		})
	}
}
