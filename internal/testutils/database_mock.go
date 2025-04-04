package testutils

import (
	"fmt"
	"testing"

	"github.com/kiasaty/spendings-tracker/models"
)

// MockDatabaseClient implements database.DatabaseClient interface
type MockDatabaseClient struct {
	spendings           map[int]*models.Spending
	tags                map[string]*models.Tag
	shouldErrorOnCreate bool
	shouldErrorOnFind   bool
}

func NewMockDatabaseClient() *MockDatabaseClient {
	return &MockDatabaseClient{
		spendings: make(map[int]*models.Spending),
		tags:      make(map[string]*models.Tag),
	}
}

type MockDatabaseClientConfig struct {
	InitialSpendings map[int]*models.Spending
	InitialTags      map[string]*models.Tag
}

func NewMockDatabaseClientWithConfig(config MockDatabaseClientConfig) *MockDatabaseClient {
	return &MockDatabaseClient{
		spendings: config.InitialSpendings,
		tags:      config.InitialTags,
	}
}

func (m *MockDatabaseClient) Migrate() {}

func (m *MockDatabaseClient) CreateTag(tag *models.Tag) (*models.Tag, error) {
	m.tags[tag.Name] = tag
	return tag, nil
}

func (m *MockDatabaseClient) FindTagByName(name string) (*models.Tag, error) {
	if tag, exists := m.tags[name]; exists {
		return tag, nil
	}
	return nil, nil
}

func (m *MockDatabaseClient) CreateSpending(spending *models.Spending) (*models.Spending, error) {
	if m.shouldErrorOnCreate {
		return nil, fmt.Errorf("mock error on create")
	}
	m.spendings[spending.MessageId] = spending
	return spending, nil
}

func (m *MockDatabaseClient) FindSpendingByMessageId(messageID int) (*models.Spending, error) {
	if spending, exists := m.spendings[messageID]; exists {
		return spending, nil
	}
	return nil, nil
}

func (m *MockDatabaseClient) UpdateSpending(spending *models.Spending) error {
	m.spendings[spending.MessageId] = spending
	return nil
}

func (m *MockDatabaseClient) SyncSpendingTags(spending *models.Spending, tags *[]models.Tag) error {
	spending.Tags = *tags
	return nil
}

func (m *MockDatabaseClient) VerifySpending(t *testing.T, messageID int, expectedCost float64) {
	spending, _ := m.FindSpendingByMessageId(messageID)
	if spending == nil {
		t.Errorf("Expected spending with message ID %d to exist", messageID)
	}
	if spending.Cost != expectedCost {
		t.Errorf("Expected cost %v, got %v", expectedCost, spending.Cost)
	}
}

func (m *MockDatabaseClient) GetSpendings() map[int]*models.Spending {
	return m.spendings
}

func (m *MockDatabaseClient) GetTags() map[string]*models.Tag {
	return m.tags
}

func (m *MockDatabaseClient) SetErrorOnCreate(shouldError bool) {
	m.shouldErrorOnCreate = shouldError
}

func (m *MockDatabaseClient) Reset() {
	m.spendings = make(map[int]*models.Spending)
	m.tags = make(map[string]*models.Tag)
}
