package extractors

import (
	"regexp"
	"strconv"
)

func ExtractHashtags(text string) []string {
	pattern := `#(\w+)`

	regex := regexp.MustCompile(pattern)

	matches := regex.FindAllString(text, -1)

	hashtags := make([]string, len(matches))
	for i, match := range matches {
		hashtags[i] = match[1:] // Remove the "#" symbol
	}

	return hashtags
}

func ExtractPrices(text string) []float64 {
	pattern := `\b\d+(\.\d+)?\b`

	regex := regexp.MustCompile(pattern)

	matches := regex.FindAllString(text, -1)

	var prices []float64

	for _, match := range matches {
		price, err := strconv.ParseFloat(match, 64)

		if err == nil {
			prices = append(prices, price)
		}
	}

	return prices
}

func ExtractDates(text string) []string {
	pattern := `\d{4}-\d{2}-\d{2}`

	regex := regexp.MustCompile(pattern)

	matches := regex.FindAllString(text, -1)

	return matches
}
