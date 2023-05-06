package pkg

import (
	"io"
	"strings"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Make an HTTP request
func FetchURLContent(url string) (string, error) {
	resp , err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err !=nil {
		return "", err
	}

	return string(body), nil
}

// Parse the HTML content and return the content
func ExtractLinksAndInfo(content string) ([]string, []interface{}) {
	links := []string{}
	// Initialize to empty
	info := []interface{}{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return links, info
	}

	// Extract links
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})

	// Extract additional info, e.g. headings, etc.
	doc.Find("h1,h2,h3,h4,h5,h6").Each(func(i int, s *goquery.Selection) {
		info = append(info, s.Text())
	})

	return links, info
}