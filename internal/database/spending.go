package database

import (
	"fmt"
	"time"

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
	err := c.DB.Where("message_id = ?", messageID).First(&spending).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find spending: %w", err)
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

func (c *Client) GetSpendingsByDateRange(startDate, endDate time.Time) ([]models.Spending, error) {
	var spendings []models.Spending
	err := c.DB.Preload("Tags").Where("spent_at BETWEEN ? AND ?", startDate, endDate).Find(&spendings).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get spendings by date range: %w", err)
	}
	return spendings, nil
}
