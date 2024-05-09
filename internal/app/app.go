package app

import (
	"fmt"
	"os"

	"github.com/kiasaty/spendings-tracker/internal/database"
)

type App struct {
	DB database.DatabaseClient
}

func NewApp(databaseClient database.DatabaseClient) App {
	return App{
		DB: databaseClient,
	}
}

func (app *App) HandleCommand() {
	if len(os.Args) < 2 {
		fmt.Println("List of existing commands.")
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
		fmt.Println("List of existing commands.")
		os.Exit(1)
	}
}
