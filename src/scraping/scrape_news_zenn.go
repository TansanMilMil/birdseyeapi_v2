package scraping

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

// ScrapeNewsByZenn implements scraping for Zenn
type ScrapeNewsByZenn struct {
	summarizer Summarizer
}

// NewScrapeNewsByZenn creates a new Zenn scraper
func NewScrapeNewsByZenn(summarizer Summarizer) *ScrapeNewsByZenn {
	return &ScrapeNewsByZenn{
		summarizer: summarizer,
	}
}

// GetSourceBy returns the source identifier
func (s *ScrapeNewsByZenn) GetSourceBy() string {
	return "Zenn"
}

// ExtractNews extracts news from Zenn
func (s *ScrapeNewsByZenn) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer // Capture summarizer for use in closure

	// URL to scrape
	url := "https://zenn.dev/topics/trending"

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Zenn: %v", err)
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
	doc.Find("article.ArticleCard_card__HlL_J").Each(func(i int, selection *goquery.Selection) {
		// Extract article URL
		articleURL := ""
		titleElement := selection.Find("h3.ArticleCard_title__TCnJm a")
		if href, exists := titleElement.Attr("href"); exists {
			articleURL = "https://zenn.dev" + href
		}

		// Extract title
		title := strings.TrimSpace(titleElement.Text())

		// Extract description
		description := strings.TrimSpace(selection.Find("div.ArticleCard_meta__ccEO7").Text())

		// Create a news item if we have at least a title and URL
		if title != "" && articleURL != "" {
			// Create a news item
			newsItem := models.News{
				Title:           title,
				Description:     description,
				SourceBy:        "Zenn",
				ScrapedUrl:      url,
				ScrapedDateTime: time.Now(),
				ArticleUrl:      articleURL,
				ArticleImageUrl: "", // Zenn articles don't typically have preview images in the listing
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