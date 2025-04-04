package app

import (
	"fmt"

	"github.com/kiasaty/spendings-tracker/models"
)

func (app *App) StoreSpending(spending *models.Spending) (*models.Spending, error) {
	spending, err := app.DB.CreateSpending(spending)
	if err != nil {
		return nil, fmt.Errorf("failed to store spending: %w", err)
	}
	return spending, nil
}

func (app *App) FindSpendingByMessageId(messageID int) (*models.Spending, error) {
	spending, err := app.DB.FindSpendingByMessageId(messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to find spending: %w", err)
	}
	return spending, nil
}

func (app *App) UpdateSpending(spending *models.Spending) (*models.Spending, error) {
	err := app.DB.UpdateSpending(spending)
	if err != nil {
		return nil, fmt.Errorf("failed to update spending: %w", err)
	}
	return spending, nil
}

func (app *App) SyncSpendingTags(spending *models.Spending, tags *[]models.Tag) error {
	err := app.DB.SyncSpendingTags(spending, tags)
	if err != nil {
		return fmt.Errorf("failed to sync spending tags: %w", err)
	}
	return nil
}
