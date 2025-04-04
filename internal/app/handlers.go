package app

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasaty/spendings-tracker/models"
	"github.com/kiasaty/spendings-tracker/pkg/extractors"
)

func (app *App) FetchUpdates() {
	updates, err := app.Bot.GetUpdates()
	if err != nil {
		fmt.Printf("Error fetching updates: %v\n", err)
		return
	}

	for _, update := range updates {
		app.handleUpdate(&update)
	}
}

func (app *App) handleUpdate(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	// Extract price, skip if not found
	price, err := extractors.ExtractPrice(update.Message.Text)
	if err != nil {
		return
	}

	// Extract date or use current time
	date, err := extractors.ExtractDate(update.Message.Text)
	if err != nil {
		date = time.Now()
	}

	// Extract tags
	tags := extractors.ExtractHashtags(update.Message.Text)
	var tagModels []models.Tag
	for _, tagName := range tags {
		tag, err := app.FindTagByName(tagName)
		if err != nil {
			fmt.Printf("Error finding tag: %v\n", err)
			continue
		}
		if tag == nil {
			tag, err = app.StoreTag(&models.Tag{
				Name: tagName,
			})
			if err != nil {
				fmt.Printf("Error storing tag: %v\n", err)
				continue
			}
		}
		tagModels = append(tagModels, *tag)
	}

	// Check if spending already exists
	spending, err := app.FindSpendingByMessageId(update.Message.MessageID)
	if err != nil {
		fmt.Printf("Error finding spending: %v\n", err)
		return
	}

	if spending == nil {
		// Create new spending
		spending, err = app.StoreSpending(&models.Spending{
			ChatId:      update.Message.Chat.ID,
			MessageId:   update.Message.MessageID,
			Cost:        price,
			Description: update.Message.Text,
			SpentAt:     date,
		})
		if err != nil {
			fmt.Printf("Error storing spending: %v\n", err)
			return
		}
	} else {
		// Update existing spending
		spending.Cost = price
		spending.Description = update.Message.Text
		spending.SpentAt = date
		spending, err = app.UpdateSpending(spending)
		if err != nil {
			fmt.Printf("Error updating spending: %v\n", err)
			return
		}
	}

	// Sync tags
	err = app.SyncSpendingTags(spending, &tagModels)
	if err != nil {
		fmt.Printf("Error syncing tags: %v\n", err)
	}
}
