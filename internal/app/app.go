package app

import (
	"fmt"
	"os"

	"github.com/kiasaty/spendings-tracker/internal/database"
	"github.com/kiasaty/spendings-tracker/pkg/telegram"
)

type App struct {
	DB  database.DatabaseClient
	Bot *telegram.TelegramBot
}

func NewApp(databaseClient database.DatabaseClient, bot *telegram.TelegramBot) (*App, error) {
	return &App{
		DB:  databaseClient,
		Bot: bot,
	}, nil
}

func (app *App) HandleCommand() {
	if len(os.Args) < 2 {
		fmt.Println("List of existing commands:")
		fmt.Println("  fetch-updates - Fetch and process new messages from Telegram")
		fmt.Println("  migrate-database - Set up the database schema")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "fetch-updates":
		app.FetchUpdates()
	case "migrate-database":
		app.DB.Migrate()
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("List of existing commands:")
		fmt.Println("  fetch-updates - Fetch and process new messages from Telegram")
		fmt.Println("  migrate-database - Set up the database schema")
		os.Exit(1)
	}
}
