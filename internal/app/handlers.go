package app

import (
	"fmt"
	"sort"
	"strings"
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

	// Handle commands
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "report":
			app.handleReportCommand(update.Message, false)
			return
		case "report_last_month":
			app.handleReportCommand(update.Message, true)
			return
		}
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

func (app *App) handleReportCommand(message *tgbotapi.Message, isLastMonth bool) {
	var startDate, endDate time.Time
	now := time.Now()

	if isLastMonth {
		// Last month's range
		startDate = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
	} else {
		// Current month's range
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	}

	// Get spendings for the period
	spendings, err := app.DB.GetSpendingsByDateRange(startDate, endDate)
	if err != nil {
		app.Bot.SendMessage(message.Chat.ID, "Failed to generate report")
		return
	}

	// Calculate totals by tag
	tagTotals := make(map[string]float64)
	var total float64

	for _, spending := range spendings {
		total += spending.Cost
		if len(spending.Tags) == 0 {
			tagTotals["no_tag"] += spending.Cost
			continue
		}
		for _, tag := range spending.Tags {
			tagTotals[tag.Name] += spending.Cost
		}
	}

	// Format the report
	var report strings.Builder
	period := "current month"
	if isLastMonth {
		period = "last month"
	}
	report.WriteString(fmt.Sprintf("Spending report for %s:\n\n", period))

	// Sort tags for consistent output
	var tags []string
	for tag := range tagTotals {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	// Add tag totals
	for _, tag := range tags {
		report.WriteString(fmt.Sprintf("%s: %.2f\n", tag, tagTotals[tag]))
	}

	// Add total
	report.WriteString(fmt.Sprintf("\nTotal: %.2f", total))

	// Send the report
	app.Bot.SendMessage(message.Chat.ID, report.String())
}
