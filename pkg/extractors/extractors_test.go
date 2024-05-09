package extractors_test

import (
	"testing"

	"github.com/kiasaty/spendings-tracker/pkg/extractors"
)

func TestExtractHashtags(t *testing.T) {
	text := "This is an #example text with #sample hashtags"

	hashtags := extractors.ExtractHashtags(text)

	if len(hashtags) != 2 {
		t.Fatal(`something went fucking wrong`)
	}
}
