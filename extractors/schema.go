package extractors

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type (
	ItemProperty struct {
		Name        string
		Content     string
		ContentText *Text
		HTML        string
	}
)

// GetAuthor retrieves the author name from a selection with itemprop=author
func GetAuthor(s *goquery.Selection) *Text {
	return GetItemProperty(s, "author").ContentText
}

// GetHeadline retrieves the healine content from a selection with itemprop=headline
func GetHeadline(s *goquery.Selection) *Text {
	return GetItemProperty(s, "headline").ContentText
}

// GetDescription retrieves the healine content from a selection with itemprop=headline
func GetDescription(s *goquery.Selection) *Text {
	return GetItemProperty(s, "description").ContentText
}

// GetItemProperty retrieves the property text
func GetItemProperty(s *goquery.Selection, property string) *ItemProperty {
	if s == nil || property == "" {
		return nil
	}
	prop := s.Find("[itemprop=" + property + "]")
	html, err := prop.Html()
	if err != nil {
		fmt.Printf("Couldn't get html for property `%s`\n", property)
		fmt.Printf("Error: %s\n", err.Error())
	}

	content := prop.AttrOr("content", "")
	if content == "" {
		content = prop.Text()
	}

	return &ItemProperty{
		Name:        property,
		Content:     html,
		ContentText: NewText(content),
	}
}
