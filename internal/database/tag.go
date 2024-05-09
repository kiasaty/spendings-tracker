package database

import "github.com/kiasaty/spendings-tracker/models"

func (c *Client) CreateTag(tag *models.Tag) (*models.Tag, error) {
	result := c.DB.Create(&tag)

	if result.Error != nil {
		return nil, result.Error
	}

	return tag, nil
}

func (c *Client) FindTagByName(name string) (*models.Tag, error) {
	tag := &models.Tag{}

	result := c.DB.Where("name = ?", name).First(tag)

	if result.Error != nil {
		return nil, result.Error
	}

	return tag, nil
}
