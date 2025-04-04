package database

import (
	"os"

	"github.com/kiasaty/spendings-tracker/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseClient interface {
	Migrate()

	CreateTag(*models.Tag) (*models.Tag, error)
	FindTagByName(name string) (*models.Tag, error)

	CreateSpending(*models.Spending) (*models.Spending, error)
	FindSpendingByMessageId(messageID int) (*models.Spending, error)
	UpdateSpending(spending *models.Spending) error
	SyncSpendingTags(*models.Spending, *[]models.Tag) error
}

type Client struct {
	DB *gorm.DB
}

func NewDatabaseClient() (DatabaseClient, error) {
	db, err := gorm.Open(
		sqlite.Open(os.Getenv("DATABASE_URL")),
		&gorm.Config{},
	)

	if err != nil {
		return nil, err
	}

	client := &Client{
		DB: db,
	}

	return client, nil
}

func (c *Client) Migrate() {
	c.DB.AutoMigrate(&models.Tag{})
	c.DB.AutoMigrate(&models.Spending{})
}
