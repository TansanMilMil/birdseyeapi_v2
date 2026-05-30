package reaction

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/env"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
)

const SOURCE_URL = "https://b.hatena.ne.jp/entry/s/"

const (
	// pageLoadTimeout bounds how long a single navigation may hang. Without
	// this a stuck page keeps the session (and its Firefox memory) alive
	// indefinitely, which is what drove the host OOM.
	pageLoadTimeout = 30 * time.Second
	// implicitWait is the per-element lookup timeout.
	implicitWait = 5 * time.Second
)

var SeleniumUrl = env.GetEnv("SELENIUM_URL", "")

// NewFirefoxDriver creates a single remote Firefox session with headless mode
// and timeouts configured. The caller owns the returned driver's lifecycle and
// must call driver.Quit() (checking the error) when done. Reusing one driver
// across many articles avoids repeatedly spawning/leaking Firefox processes.
func NewFirefoxDriver() (selenium.WebDriver, error) {
	remoteURL, err := url.Parse(SeleniumUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid selenium URL: %v", err)
	}

	caps := selenium.Capabilities{"browserName": "firefox"}
	caps.AddFirefox(firefox.Capabilities{
		Args: []string{"-headless"},
	})

	driver, err := selenium.NewRemote(caps, remoteURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create selenium session: %v", err)
	}

	if err := driver.SetPageLoadTimeout(pageLoadTimeout); err != nil {
		// Best effort: log via the returned error path by quitting and failing.
		_ = driver.Quit()
		return nil, fmt.Errorf("failed to set page load timeout: %v", err)
	}
	if err := driver.SetImplicitWaitTimeout(implicitWait); err != nil {
		_ = driver.Quit()
		return nil, fmt.Errorf("failed to set implicit wait timeout: %v", err)
	}

	return driver, nil
}

type ScrapeReactionsByHatena struct{}

func NewScrapeReactionsByHatena() *ScrapeReactionsByHatena {
	return &ScrapeReactionsByHatena{}
}

func (s *ScrapeReactionsByHatena) GetSourceBy() string {
	return "Hatena"
}

func (s *ScrapeReactionsByHatena) ExtractReactions(driver selenium.WebDriver, newsId uint, articleURL string, title string) ([]models.NewsReaction, error) {
	var reactions []models.NewsReaction

	// Clean the URL (remove protocol)
	cleanURL := articleURL
	cleanURL = strings.Replace(cleanURL, "http://", "", 1)
	cleanURL = strings.Replace(cleanURL, "https://", "", 1)
	hatenaURL := SOURCE_URL + cleanURL

	if err := driver.Get(hatenaURL); err != nil {
		return nil, fmt.Errorf("failed to navigate to URL: %v", err)
	}

	// Wait for page to load
	time.Sleep(1 * time.Second)

	// Find comments using CSS selector
	articleSelector := "#container > div > div.entry-contents > div.entry-main > div.entry-comments > div > div.bookmarks-sort-panels.js-bookmarks-sort-panels > div.is-active.bookmarks-sort-panel.js-bookmarks-sort-panel > div > div > div.entry-comment-contents-main > .entry-comments-contents-body > .js-bookmark-comment"
	articleElements, err := driver.FindElements(selenium.ByCSSSelector, articleSelector)
	if err != nil {
		fmt.Printf("error finding elements: %v\n", err)
	}

	// Process each comment
	for _, article := range articleElements {
		text, err := article.Text()
		if err != nil {
			fmt.Printf("failed to get text: %v\n", err)
			continue
		}

		if text == "" || strings.TrimSpace(text) == "" || text == title {
			continue
		}

		// Create NewsReaction object
		reaction := models.NewsReaction{
			NewsID:          newsId,
			Author:          "hatena user",
			Comment:         strings.TrimSpace(text),
			ScrapedDateTime: time.Now().UTC(),
			CommentUrl:      hatenaURL,
		}

		reactions = append(reactions, reaction)
	}

	return reactions, nil
}
