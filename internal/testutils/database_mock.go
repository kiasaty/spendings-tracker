package testutils

import (
	"fmt"
	"testing"
	"time"

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

func (m *MockDatabaseClient) VerifySpending(t *testing.T, spending *models.Spending, expectedCost float64, expectedDate time.Time) {
	if spending == nil {
		t.Errorf("Expected spending to exist")
		return
	}
	if spending.Cost != expectedCost {
		t.Errorf("Expected cost %v, got %v", expectedCost, spending.Cost)
	}
	// Truncate both dates to seconds for comparison
	expectedDateTruncated := expectedDate.Truncate(time.Second)
	actualDateTruncated := spending.SpentAt.Truncate(time.Second)
	if !actualDateTruncated.Equal(expectedDateTruncated) {
		t.Errorf("Expected date %v, got %v", expectedDateTruncated, actualDateTruncated)
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

func (m *MockDatabaseClient) FindTagsBySpendingId(spendingID uint) ([]models.Tag, error) {
	for _, spending := range m.spendings {
		if spending.ID == spendingID {
			return spending.Tags, nil
		}
	}
	return nil, nil
}

func (m *MockDatabaseClient) VerifySpendingTags(t *testing.T, spending *models.Spending, expectedTags []string) {
	if spending == nil {
		t.Errorf("Expected spending to exist")
		return
	}

	tags, _ := m.FindTagsBySpendingId(spending.ID)
	if len(tags) != len(expectedTags) {
		t.Errorf("Expected %d tags, got %d", len(expectedTags), len(tags))
		return
	}

	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag.Name] = true
	}

	for _, expectedTag := range expectedTags {
		if !tagMap[expectedTag] {
			t.Errorf("Expected tag %s not found", expectedTag)
		}
	}
}

func (m *MockDatabaseClient) GetSpendingsByDateRange(startDate, endDate time.Time) ([]models.Spending, error) {
	var result []models.Spending
	for _, spending := range m.spendings {
		if !spending.SpentAt.Before(startDate) && !spending.SpentAt.After(endDate) {
			result = append(result, *spending)
		}
	}
	return result, nil
}
