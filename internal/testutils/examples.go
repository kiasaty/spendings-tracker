package testutils

import (
	"fmt"

	"github.com/kiasaty/spendings-tracker/models"
)

// ExampleMockDatabaseClient demonstrates how to use the mock database client
func ExampleMockDatabaseClient() {
	mockDB := NewMockDatabaseClient()

	// Create a spending
	spending := &models.Spending{
		MessageId:   1,
		ChatId:      123456789,
		Cost:        15.50,
		Description: "Lunch",
	}

	// Store the spending
	_, err := mockDB.CreateSpending(spending)
	if err != nil {
		fmt.Println("Error creating spending:", err)
		return
	}

	// Retrieve the spending
	retrieved, _ := mockDB.FindSpendingByMessageId(1)
	fmt.Printf("Retrieved spending cost: %.2f\n", retrieved.Cost)
	// Output: Retrieved spending cost: 15.50
}
