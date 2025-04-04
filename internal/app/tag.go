package app

import (
	"fmt"

	"github.com/kiasaty/spendings-tracker/models"
)

func (app *App) StoreTag(tag *models.Tag) (*models.Tag, error) {
	tag, err := app.DB.CreateTag(tag)
	if err != nil {
		return nil, fmt.Errorf("failed to store tag: %w", err)
	}
	return tag, nil
}

func (app *App) FindTagByName(name string) (*models.Tag, error) {
	tag, err := app.DB.FindTagByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}
	return tag, nil
}
