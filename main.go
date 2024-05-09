package main

import (
	"github.com/joho/godotenv"
	"github.com/kiasaty/spendings-tracker/internal/app"
	"github.com/kiasaty/spendings-tracker/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Loading .env file failed!")
	}

	databaseClient, err := database.NewDatabaseClient()

	if err != nil {
		panic("Setting up database client failed!")
	}

	app := app.NewApp(databaseClient)

	app.HandleCommand()
}
