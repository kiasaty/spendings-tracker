package app

import "github.com/kiasaty/spendings-tracker/models"

func (app *App) StoreSpending(spending *models.Spending) *models.Spending {
	spending, err := app.DB.CreateSpending(spending)

	if err != nil {
		panic(err.Error())
	}

	return spending
}

func (app *App) FindSpendingByMessageId(messageID int) (spending *models.Spending) {
	return app.DB.FindSpendingByMessageId(messageID)
}

func (app *App) UpdateSpending(spending *models.Spending) *models.Spending {
	err := app.DB.UpdateSpending(spending)

	if err != nil {
		panic(err.Error())
	}

	return spending
}

func (app *App) SyncSpendingTags(spending *models.Spending, tags *[]models.Tag) {
	app.DB.SyncSpendingTags(spending, tags)
}
