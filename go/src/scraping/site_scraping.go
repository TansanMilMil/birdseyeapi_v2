package scraping

import (
	"fmt"
	"net/url"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/ai"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping/news"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping/reaction"
	"github.com/tebeka/selenium"
)

type SiteScraping struct {
	scrapers         []news.ScrapingNews
	reactionScrapers []reaction.ScrapingReaction
}

func NewSiteScraping() *SiteScraping {
	summarizer := ai.NewOpenAISummarizer()

	return &SiteScraping{
		scrapers: []news.ScrapingNews{
			news.NewScrapeNewsByCloudWatchImpress(summarizer),
			news.NewScrapeNewsByHatena(summarizer),
			news.NewScrapeNewsByZenn(summarizer),
			news.NewScrapeNewsByZDNet(summarizer),
		},
		reactionScrapers: []reaction.ScrapingReaction{
			reaction.NewScrapeReactionsByHatena(),
		},
	}
}

func (s *SiteScraping) ScrapeNews() ([]models.News, error) {
	allNews := []models.News{}

	for _, scraper := range s.scrapers {
		fmt.Print(scraper.GetSourceBy() + ": scraping")
		news, err := safeExtractNews(scraper)
		if err != nil {
			fmt.Printf(" -> Error scraping from %s: %v\n", scraper.GetSourceBy(), err)
			continue
		}

		fmt.Printf(" -> scraped article: %s: %d\n", scraper.GetSourceBy(), len(news))
		allNews = append(allNews, news...)
	}

	return allNews, nil
}

// NewReactionDriver creates a single Selenium session to be shared across all
// reaction scrapes in one run. The caller owns its lifecycle and must Quit it.
func (s *SiteScraping) NewReactionDriver() (selenium.WebDriver, error) {
	return reaction.NewFirefoxDriver()
}

func (s *SiteScraping) ScrapeReactions(driver selenium.WebDriver, news models.News) ([]models.NewsReaction, error) {
	allReactions := []models.NewsReaction{}

	_, err := url.Parse(news.ArticleUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid article URL: %v", err)
	}

	for _, scraper := range s.reactionScrapers {
		reactions, err := safeExtractReactions(scraper, driver, news.ID, news.ArticleUrl, news.Title)
		if err != nil {
			fmt.Printf("Error scraping reactions from %s: %v\n", scraper.GetSourceBy(), err)
			continue
		}

		fmt.Printf(" -> scraped reactions: %s: %d\n", scraper.GetSourceBy(), len(reactions))
		allReactions = append(allReactions, reactions...)
	}

	return allReactions, nil
}

// safeExtractNews runs a single news scraper, converting any panic into an
// error so that a runtime failure in one site does not abort the scraping of
// the remaining sites.
func safeExtractNews(scraper news.ScrapingNews) (result []models.News, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("panic while scraping %s: %v", scraper.GetSourceBy(), r)
		}
	}()
	return scraper.ExtractNews()
}

// safeExtractReactions runs a single reaction scraper, converting any panic
// into an error so one site's runtime failure does not abort the others.
func safeExtractReactions(scraper reaction.ScrapingReaction, driver selenium.WebDriver, newsID uint, articleURL, title string) (result []models.NewsReaction, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("panic while scraping reactions from %s: %v", scraper.GetSourceBy(), r)
		}
	}()
	return scraper.ExtractReactions(driver, newsID, articleURL, title)
}
