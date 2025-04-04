package database

import (
	"fmt"

	"github.com/kiasaty/spendings-tracker/models"
	"gorm.io/gorm"
)

func (c *Client) CreateTag(tag *models.Tag) (*models.Tag, error) {
	result := c.DB.Create(&tag)

	if result.Error != nil {
		return nil, result.Error
	}

	return tag, nil
}

func (c *Client) FindTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := c.DB.Where("name = ?", name).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tag: %w", err)
	}
	return &tag, nil
}
