package main

import (
	"github.com/joho/godotenv"
	"github.com/kiasaty/spendings-tracker/internal/app"
	"github.com/kiasaty/spendings-tracker/internal/database"
	"github.com/kiasaty/spendings-tracker/pkg/telegram"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Loading .env file failed!")
	}

	databaseClient, err := database.NewDatabaseClient()
	if err != nil {
		panic("Setting up database client failed!")
	}

	bot, err := telegram.NewTelegramBot()
	if err != nil {
		panic("failed to create telegram bot: " + err.Error())
	}

	app, err := app.NewApp(databaseClient, bot)
	if err != nil {
		panic("Setting up app failed: " + err.Error())
	}

	app.HandleCommand()
}
