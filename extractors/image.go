package extractors

import (
	"github.com/PuerkitoBio/goquery"
)

type (
	Image struct {
		// Data   []byte
		SrcURL string
	}
)

// GetImage Retrieves an image instance from a given selector
func GetImage(s *goquery.Selection, selector string) *Image {
	if s == nil {
		return nil
	}

	return &Image{
		SrcURL: s.Find(selector).AttrOr("src", ""),
	}
}
