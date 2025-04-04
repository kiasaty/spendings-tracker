package database

import (
	"errors"

	"github.com/kiasaty/spendings-tracker/models"
	"gorm.io/gorm"
)

func (c *Client) CreateSpending(spending *models.Spending) (*models.Spending, error) {
	result := c.DB.Create(&spending)

	if result.Error != nil {
		return nil, result.Error
	}

	return spending, nil
}

func (c *Client) FindSpendingByMessageId(messageID int) (*models.Spending, error) {
	var spending models.Spending
	result := c.DB.Where("message_id = ?", messageID).First(&spending)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &spending, nil
}

func (c *Client) UpdateSpending(spending *models.Spending) error {
	result := c.DB.Save(&spending)

	return result.Error
}

func (c *Client) SyncSpendingTags(spending *models.Spending, tags *[]models.Tag) error {
	return c.DB.Model(spending).Association("Tags").Replace(tags)
}
