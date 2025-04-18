package scraping

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/internal/models"
)

// ScrapeNewsByCloudWatchImpress implements scraping for CloudWatch by Impress
type ScrapeNewsByCloudWatchImpress struct {
	summarizer Summarizer
}

// NewScrapeNewsByCloudWatchImpress creates a new CloudWatch by Impress scraper
func NewScrapeNewsByCloudWatchImpress(summarizer Summarizer) *ScrapeNewsByCloudWatchImpress {
	return &ScrapeNewsByCloudWatchImpress{
		summarizer: summarizer,
	}
}

// GetSourceBy returns the source identifier
func (s *ScrapeNewsByCloudWatchImpress) GetSourceBy() string {
	return "CloudWatch by Impress"
}

// ExtractNews extracts news from CloudWatch by Impress
func (s *ScrapeNewsByCloudWatchImpress) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer // Capture summarizer for use in closure

	// URL to scrape
	url := "https://cloud.watch.impress.co.jp/"

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CloudWatch: %v", err)
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
	doc.Find("div.article").Each(func(i int, selection *goquery.Selection) {
		// Extract article URL
		articleURL := ""
		titleElement := selection.Find("h3 a")
		if href, exists := titleElement.Attr("href"); exists {
			articleURL = href
			if !strings.HasPrefix(articleURL, "http") {
				articleURL = "https://cloud.watch.impress.co.jp" + articleURL
			}
		}

		// Extract title
		title := strings.TrimSpace(titleElement.Text())

		// Extract description
		description := strings.TrimSpace(selection.Find("p.description").Text())

		// Extract image URL
		imageURL := ""
		if src, exists := selection.Find("img").Attr("src"); exists {
			imageURL = src
			if !strings.HasPrefix(imageURL, "http") {
				imageURL = "https://cloud.watch.impress.co.jp" + imageURL
			}
		}

		// Create a news item if we have at least a title and URL
		if title != "" && articleURL != "" {
			// Create a news item
			newsItem := models.News{
				Title:           title,
				Description:     description,
				SourceBy:        "CloudWatch by Impress",
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