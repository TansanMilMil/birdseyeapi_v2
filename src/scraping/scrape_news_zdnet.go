package scraping

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

// ScrapeNewsByZDNet implements scraping for ZDNet Japan
type ScrapeNewsByZDNet struct {
	summarizer Summarizer
}

// NewScrapeNewsByZDNet creates a new ZDNet Japan scraper
func NewScrapeNewsByZDNet(summarizer Summarizer) *ScrapeNewsByZDNet {
	return &ScrapeNewsByZDNet{
		summarizer: summarizer,
	}
}

// GetSourceBy returns the source identifier
func (s *ScrapeNewsByZDNet) GetSourceBy() string {
	return "ZDNet Japan"
}

// ExtractNews extracts news from ZDNet Japan
func (s *ScrapeNewsByZDNet) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer // Capture summarizer for use in closure

	// URL to scrape
	url := "https://japan.zdnet.com/topics/"

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ZDNet: %v", err)
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
	doc.Find("div.article-item").Each(func(i int, selection *goquery.Selection) {
		// Extract article URL
		articleURL := ""
		titleElement := selection.Find("h3.title a")
		if href, exists := titleElement.Attr("href"); exists {
			articleURL = href
			if !strings.HasPrefix(articleURL, "http") {
				articleURL = "https://japan.zdnet.com" + articleURL
			}
		}

		// Extract title
		title := strings.TrimSpace(titleElement.Text())

		// Extract description
		description := strings.TrimSpace(selection.Find("p.abstract").Text())

		// Extract image URL
		imageURL := ""
		if src, exists := selection.Find("div.image img").Attr("src"); exists {
			imageURL = src
		}

		// Create a news item if we have at least a title and URL
		if title != "" && articleURL != "" {
			// Create a news item
			newsItem := models.News{
				Title:           title,
				Description:     description,
				SourceBy:        "ZDNet Japan",
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