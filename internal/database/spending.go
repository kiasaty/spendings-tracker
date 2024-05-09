package database

import "github.com/kiasaty/spendings-tracker/models"

func (c *Client) CreateSpending(spending *models.Spending) (*models.Spending, error) {
	result := c.DB.Create(&spending)

	if result.Error != nil {
		return nil, result.Error
	}

	return spending, nil
}

func (c *Client) FindSpendingByMessageId(messageID int) (spending *models.Spending) {
	c.DB.Where("message_id = ?", messageID).First(&spending)

	return spending
}

func (c *Client) UpdateSpending(spending *models.Spending) error {
	result := c.DB.Save(&spending)

	return result.Error
}

func (c *Client) SyncSpendingTags(spending *models.Spending, tags *[]models.Tag) error {
	return c.DB.Model(*spending).Association("Tags").Replace(tags)
}
