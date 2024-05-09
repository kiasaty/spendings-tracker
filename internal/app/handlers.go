package app

import (
	"github.com/kiasaty/spendings-tracker/models"
	"github.com/kiasaty/spendings-tracker/pkg/extractors"
	"github.com/kiasaty/spendings-tracker/pkg/telegram"
)

func (app *App) FetchUpdates() {
	updates := telegram.GetUpdates()

	for _, update := range updates {
		app.handleUpdate(&update)
	}
}

func (app *App) handleUpdate(update *telegram.Update) {
	prices := extractors.ExtractPrices(update.Message.Text)

	if len(prices) != 1 {
		return
	}

	var tags []models.Tag

	tagsNames := extractors.ExtractHashtags(update.Message.Text)

	for _, tagName := range tagsNames {
		tag := app.FindTagByName(tagName)

		if tag == nil {
			tag = app.StoreTag(&models.Tag{
				Name: tagName,
			})
		}

		tags = append(tags, *tag)
	}

	spending := app.FindSpendingByMessageId(update.Message.MessageID)

	if spending == nil {
		spending = app.StoreSpending(&models.Spending{
			ChatId:      update.Message.Chat.ID,
			MessageId:   update.Message.MessageID,
			Cost:        prices[0],
			Description: update.Message.Text,
		})
	} else {
		spending = app.UpdateSpending(&models.Spending{
			Cost:        prices[0],
			Description: update.Message.Text,
		})
	}

	app.SyncSpendingTags(spending, &tags)
}
