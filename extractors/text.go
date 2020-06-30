package extractors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type (
	Text struct {
		*string
		Trimmed string
		Length  int
	}
)

func NewText(s string) *Text {
	if s == "" {
		return nil
	}
	return &Text{
		string:  &s,
		Length:  len(s),
		Trimmed: strings.TrimSpace(s),
	}
}

func (t *Text) Split(sep string) []Text {
	result := make([]Text, 0)

	parts := strings.Split(t.Trimmed, sep)
	if len(parts) > 0 {
		for _, part := range parts {
			result = append(result, *NewText(part))
		}
	}

	return result
}

// GetText Retrieves text from the given document already trimmed
func GetText(s *goquery.Selection, selector string) *Text {
	if s == nil {
		return nil
	}
	return NewText(s.Find(selector).Text())
}

// AsInt Retrieve Text as integer number
func (t *Text) AsInt() int {
	if t == nil || t.Length == 0 {
		return 0
	}

	result, err := strconv.Atoi(t.Trimmed)
	if err != nil {
		fmt.Printf("Could not convert text: %s to integer\n", t.Trimmed)
		fmt.Printf("Error: %s\n", err.Error())
	}

	return result
}
