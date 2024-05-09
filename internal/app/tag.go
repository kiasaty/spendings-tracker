package app

import "github.com/kiasaty/spendings-tracker/models"

func (app *App) StoreTag(tag *models.Tag) *models.Tag {
	tag, err := app.DB.CreateTag(tag)

	if err != nil {
		panic(err.Error())
	}

	return tag
}

func (app *App) FindTagByName(name string) *models.Tag {
	tag, err := app.DB.FindTagByName(name)

	if err != nil {
		panic(err.Error())
	}

	return tag
}
