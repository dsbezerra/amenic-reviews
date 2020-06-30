package helpers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
)

// UserAgents is a list of user agents used by scraper
var UserAgents = [...]string{
	// Linus Firefox
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0",
	// Mac Firefox
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:43.0) Gecko/20100101 Firefox/43.0",
	// Mac Safari 4
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_2; de-at) AppleWebKit/531.21.8 (KHTML, like Gecko) Version/4.0.4 Safari/531.21.10",
	// Mac Safari
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
	// Windows Chrome
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.125 Safari/537.36",
	// Windows IE 10
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; WOW64; Trident/6.0)",
	// Windows IE 11
	"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
	// Windows Edge
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586",
	// Windows Firefox
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:43.0) Gecko/20100101 Firefox/43.0",
	// iPhone
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B5110e Safari/601.1",
	// iPad
	"Mozilla/5.0 (iPad; CPU OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	// Android
	"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 7 Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.76 Safari/537.36",
}

// GetRandomUserAgent retrieves a random user agent
func GetRandomUserAgent() string {
	result := ""

	// Using current time nanosecond as seed
	seed := time.Now().Nanosecond()

	// Seed the random
	rand.Seed(int64(seed))

	// Get random user-agent
	size := len(UserAgents)
	result = UserAgents[rand.Int31n(int32(size))]

	return result
}

// ResolveURL resolves a relative url to absolute
func ResolveURL(base, relative string) string {
	// Converts relative string url to URL type
	u, err := url.Parse(relative)
	if err != nil {
		log.Fatal(err)
	}

	// Convert base string url to URL type
	bu, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}

	// Resolves the relative to absolute url as string
	return bu.ResolveReference(u).String()
}

// NewDocument gets a new goquery.Document from a given website
func NewDocument(url, charset string) (*goquery.Document, error) {
	if charset != "" && charset != "utf-8" {
		return NewConvertedDocument(url, charset)
	}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// NewConvertedDocument gets a new goquery.Document from a given website coverted to utf-8
func NewConvertedDocument(url, charset string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", GetRandomUserAgent())

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		fmt.Printf(response.Status)
		// TODO
		log.Fatal(errors.New("debug"))
	}

	defer response.Body.Close()

	// Convert charset to utf-8
	utfBody, err := iconv.NewReader(response.Body, charset, "utf-8")
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
