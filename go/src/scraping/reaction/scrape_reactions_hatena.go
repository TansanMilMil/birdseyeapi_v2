package reaction

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/env"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/sclevine/agouti"
)

const SOURCE_URL = "https://b.hatena.ne.jp/entry/"

var SeleniumUrl = env.GetEnv("SELENIUM_URL", "")

type ScrapeReactionsByHatena struct{}

func NewScrapeReactionsByHatena() *ScrapeReactionsByHatena {
	return &ScrapeReactionsByHatena{}
}

func (s *ScrapeReactionsByHatena) GetSourceBy() string {
	return "Hatena"
}

func (s *ScrapeReactionsByHatena) ExtractReactions(articleURL, title string) ([]models.NewsReaction, error) {
	var reactions []models.NewsReaction

	// If SeleniumUrl is specified, use remote driver instead
	remoteURL, err := url.Parse(SeleniumUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid selenium URL: %v", err)
	}
	options := agouti.ChromeOptions(
		"args", []string{"--headless", "--disable-gpu", "--no-sandbox"},
	)
	driver := agouti.NewWebDriver(remoteURL.String(), []string{}, options)

	err = driver.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start selenium driver: %v", err)
	}
	defer driver.Stop()

	fmt.Println("selenium is ready.")

	// Clean the URL (remove protocol)
	cleanURL := articleURL
	cleanURL = strings.Replace(cleanURL, "http://", "", 1)
	cleanURL = strings.Replace(cleanURL, "https://", "", 1)
	hatenaURL := SOURCE_URL + cleanURL

	// Create a new page and navigate to the Hatena URL
	page, err := driver.NewPage()
	if err != nil {
		return nil, fmt.Errorf("failed to create new page: %v", err)
	}

	fmt.Println("selenium is requesting hatena.")
	err = page.Navigate(hatenaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to URL: %v", err)
	}

	// Wait for page to load
	time.Sleep(1 * time.Second)
	fmt.Println("request completed.")

	// Find comments using CSS selector
	articleSelector := "#container > div > div.entry-contents > div.entry-main > div.entry-comments > div > div.bookmarks-sort-panels.js-bookmarks-sort-panels > div.is-active.bookmarks-sort-panel.js-bookmarks-sort-panel > div > div > div.entry-comment-contents-main > .entry-comments-contents-body > .js-bookmark-comment"
	articleElements, err := page.All(articleSelector).Elements()
	if err != nil {
		fmt.Printf("error finding elements: %v\n", err)
	}

	fmt.Printf("articles.size(): %d\n", len(articleElements))

	// Process each comment
	for _, article := range articleElements {
		text, err := article.GetText()
		if err != nil {
			fmt.Printf("failed to get text: %v\n", err)
			continue
		}

		if text == "" || strings.TrimSpace(text) == "" || text == title {
			continue
		}

		fmt.Println("-------------------------")
		fmt.Printf("textLength: %d\n", len(text))

		// Create NewsReaction object
		reaction := models.NewsReaction{
			ReactionType: "comment",
			Source:       "hatena user",
			Count:        1,
			ScrapedAt:    time.Now().UTC(),
		}

		reactions = append(reactions, reaction)
	}

	fmt.Println("selenium quit.")
	return reactions, nil
}
