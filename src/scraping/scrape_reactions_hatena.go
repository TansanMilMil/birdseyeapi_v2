package scraping

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

type ScrapeReactionsByHatena struct{}

func NewScrapeReactionsByHatena() *ScrapeReactionsByHatena {
	return &ScrapeReactionsByHatena{}
}

func (s *ScrapeReactionsByHatena) GetSourceBy() string {
	return "Hatena"
}

func (s *ScrapeReactionsByHatena) ExtractReactions(articleURL, title string) ([]models.NewsReaction, error) {
	var reactions []models.NewsReaction

	encodedURL := url.QueryEscape(articleURL)
	hatenaURL := "https://b.hatena.ne.jp/entry/json/" + encodedURL

	resp, err := http.Get(hatenaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Hatena reactions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	bookmarkCountText := strings.TrimSpace(doc.Find("span.bookmark-count").Text())
	bookmarkCount := 0
	if bookmarkCountText != "" {
		count, err := strconv.Atoi(bookmarkCountText)
		if err == nil {
			bookmarkCount = count
		}
	}

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
