package extractors

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
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

func ExtractPrice(text string) (float64, error) {
	pattern := `\b\d+(\.\d+)?\b`

	regex := regexp.MustCompile(pattern)

	matches := regex.FindAllString(text, -1)

	for _, match := range matches {
		price, err := strconv.ParseFloat(match, 64)

		if err != nil {
			continue
		}

		return price, nil
	}

	return 0, fmt.Errorf("no price was found")
}

func ExtractDate(text string) (time.Time, error) {
	patterns := []struct {
		pattern string
		layout  string
	}{
		{`\d{4}-\d{2}-\d{2}`, "2006-01-02"},
		{`\d{2}-\d{2}-\d{4}`, "02-01-2006"},
		{`\d{2}\.\d{2}\.\d{4}`, "02.01.2006"},
		{`\d{2}/\d{2}/\d{4}`, "01/02/2006"},
	}

	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern.pattern)

		matches := regex.FindAllString(text, -1)

		for _, match := range matches {
			date, err := time.Parse(pattern.layout, match)

			if err != nil {
				continue
			}

			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("no date was found in the text")
}
