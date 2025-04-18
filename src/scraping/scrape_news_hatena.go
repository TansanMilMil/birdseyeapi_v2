package scraping

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/internal/models"
)

// ScrapeNewsByHatena implements scraping for Hatena
type ScrapeNewsByHatena struct {
	summarizer Summarizer
}

// NewScrapeNewsByHatena creates a new Hatena scraper
func NewScrapeNewsByHatena(summarizer Summarizer) *ScrapeNewsByHatena {
	return &ScrapeNewsByHatena{
		summarizer: summarizer,
	}
}

// GetSourceBy returns the source identifier
func (s *ScrapeNewsByHatena) GetSourceBy() string {
	return "Hatena"
}

// ExtractNews extracts news from Hatena
func (s *ScrapeNewsByHatena) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer // Capture summarizer for use in closure

	// URL to scrape
	url := "https://b.hatena.ne.jp/hotentry/it"

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Hatena: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Find news articles
	doc.Find("div.entrylist-item").Each(func(i int, selection *goquery.Selection) {
		// Extract article URL
		articleURL := ""
		titleElement := selection.Find("h3.entrylist-contents-title a")
		if href, exists := titleElement.Attr("href"); exists {
			articleURL = href
		}

		// Extract title
		title := strings.TrimSpace(titleElement.Text())

		// Extract description
		description := strings.TrimSpace(selection.Find("div.entrylist-contents-body").Text())

		// Extract image URL
		imageURL := ""
		if src, exists := selection.Find("img.entrylist-contents-thumb").Attr("src"); exists {
			imageURL = src
		}

		// Create a news item if we have at least a title and URL
		if title != "" && articleURL != "" {
			// Create a news item
			newsItem := models.News{
				Title:           title,
				Description:     description,
				SourceBy:        "Hatena",
				ScrapedUrl:      url,
				ScrapedDateTime: time.Now(),
				ArticleUrl:      articleURL,
				ArticleImageUrl: imageURL,
			}

			// Try to summarize the description if available
			if description != "" && summarizer != nil {
				summary, err := summarizer.Summarize(description)
				if err == nil {
					newsItem.SummarizedText = summary
				}
			}

			news = append(news, newsItem)
		}
	})

	return news, nil
}