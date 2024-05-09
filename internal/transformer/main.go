package handlers

import "github.com/kiasaty/spendings-tracker/models"

func TextToSpending(text string) models.Spending {
	var spending models.Spending

	spending.Description = text

	return spending
}
