package reaction

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
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

	// remove http:// or https:// from the URL
	encodedURL := strings.TrimPrefix(articleURL, "http://")
	encodedURL = strings.TrimPrefix(encodedURL, "https://")
	hatenaURL := "https://b.hatena.ne.jp/entry/s/" + encodedURL

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
