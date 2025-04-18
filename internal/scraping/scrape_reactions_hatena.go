package scraping

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/internal/models"
)

// ScrapeReactionsByHatena implements reaction scraping for Hatena
type ScrapeReactionsByHatena struct {}

// NewScrapeReactionsByHatena creates a new Hatena reaction scraper
func NewScrapeReactionsByHatena() *ScrapeReactionsByHatena {
	return &ScrapeReactionsByHatena{}
}

// GetSourceBy returns the source identifier
func (s *ScrapeReactionsByHatena) GetSourceBy() string {
	return "Hatena"
}

// ExtractReactions extracts reactions for a specific article from Hatena
func (s *ScrapeReactionsByHatena) ExtractReactions(articleURL, title string) ([]models.NewsReaction, error) {
	var reactions []models.NewsReaction

	// Encode the URL for use in the Hatena API
	encodedURL := url.QueryEscape(articleURL)
	hatenaURL := "https://b.hatena.ne.jp/entry/json/" + encodedURL

	// Make HTTP request
	resp, err := http.Get(hatenaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Hatena reactions: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	// Parse the response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Try to extract bookmark count
	bookmarkCountText := strings.TrimSpace(doc.Find("span.bookmark-count").Text())
	bookmarkCount := 0
	if bookmarkCountText != "" {
		count, err := strconv.Atoi(bookmarkCountText)
		if err == nil {
			bookmarkCount = count
		}
	}

	// Create a reaction item for bookmarks if count > 0
	if bookmarkCount > 0 {
		reaction := models.NewsReaction{
			ReactionType: "bookmark",
			Count:        bookmarkCount,
			Source:       "Hatena",
			ScrapedAt:    time.Now(),
		}
		reactions = append(reactions, reaction)
	}

	return reactions, nil
}