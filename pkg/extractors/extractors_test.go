package extractors_test

import (
	"reflect"
	"testing"

	"github.com/kiasaty/spendings-tracker/pkg/extractors"
)

func TestExtractHashtags(t *testing.T) {
	tests := []struct {
		testName  string
		inputText string
		expected  []string
	}{
		{
			testName:  "it extracts the hashtag from a text",
			inputText: "This is an #example text with a hashtag",
			expected:  []string{"example"},
		},
		{
			testName:  "it extracts all the hashtags in a text",
			inputText: "This is an #example text with #sample hashtags",
			expected:  []string{"example", "sample"},
		},
		{
			testName:  "it can handle text containing only a hashtag",
			inputText: "#example",
			expected:  []string{"example"},
		},
		{
			testName:  "it returns an empty list when there is no hashtag in a text",
			inputText: "This is an example text with no hashtags in it",
			expected:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			hashtags := extractors.ExtractHashtags(tt.inputText)

			if !reflect.DeepEqual(hashtags, tt.expected) {
				t.Fail()
			}
		})

	}
}

func TestExtractPrices(t *testing.T) {
	tests := []struct {
		testName      string
		inputText     string
		expectedPrice float64
		expectedError string
	}{
		{
			testName:      "it extracts the price from a text",
			inputText:     "This is an example text with prices like 2.50 in it",
			expectedPrice: 2.50,
			expectedError: "",
		},
		{
			testName:      "it returns the first found price in a text",
			inputText:     "An example text with two prices like 3.40 and 1.50 in it",
			expectedPrice: 3.40,
			expectedError: "",
		},
		{
			testName:      "it supports numbers without decimal",
			inputText:     "2 euros",
			expectedPrice: 2.00,
			expectedError: "",
		},
		{
			testName:      "it supports number with one decimal",
			inputText:     "10.1 euros",
			expectedPrice: 10.10,
			expectedError: "",
		},
		{
			testName:      "it returns no-price-found error when no price can be found in the text",
			inputText:     "this is an example text with no price in it",
			expectedPrice: 0,
			expectedError: "no price was found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			price, err := extractors.ExtractPrice(tt.inputText)

			if err != nil {
				if tt.expectedError == "" {
					t.Error(err.Error())
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%v'", tt.expectedError, err)
				}

				return
			}

			if price != tt.expectedPrice {
				t.Fatalf("Expected the price to be %.2f, but got %.2f", tt.expectedPrice, price)
			}
		})
	}
}

func TestExtractDate(t *testing.T) {
	tests := []struct {
		testName      string
		inputText     string
		expectedDate  string
		expectedError string
	}{
		{
			testName:      "it extracts a date from a text",
			inputText:     "This is an example text with a date like 2024-05-09 in it",
			expectedDate:  "2024-05-09",
			expectedError: "",
		},
		{
			testName:      "it can find dates in 2006-01-02 format",
			inputText:     "2024-05-09",
			expectedDate:  "2024-05-09",
			expectedError: "",
		},
		{
			testName:      "it can find dates in 02-01-2006 format",
			inputText:     "09-05-2024",
			expectedDate:  "2024-05-09",
			expectedError: "",
		},
		{
			testName:      "it can find dates in 02.01.2006 format",
			inputText:     "09.05.2024",
			expectedDate:  "2024-05-09",
			expectedError: "",
		},
		{
			testName:      "it can find dates in 01/02/2006 format",
			inputText:     "05/09/2024",
			expectedDate:  "2024-05-09",
			expectedError: "",
		},
		{
			testName:      "it returns no-date-found error when no date has been found",
			inputText:     "This is an example text without any date in it",
			expectedDate:  "",
			expectedError: "no date was found in the text",
		},
		{
			testName:      "it ignores invalid dates",
			inputText:     "25/09/2024 33.05.2024 33-05-0000 00-00-0000",
			expectedDate:  "",
			expectedError: "no date was found in the text",
		},
		{
			testName:      "it returns the first found valid date in a text",
			inputText:     "25/09/2024 2024-13-05 2024-05-09 and 2024-05-10",
			expectedDate:  "2024-05-09",
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			date, err := extractors.ExtractDate(tt.inputText)

			if err != nil {
				if tt.expectedError == "" {
					t.Error(err.Error())
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%v'", tt.expectedError, err)
				}

				return
			}

			formattedDate := date.Format("2006-01-02")

			if formattedDate != tt.expectedDate {
				t.Errorf("Expected the date to be %s, but got %s", tt.expectedDate, formattedDate)
			}
		})
	}
}
