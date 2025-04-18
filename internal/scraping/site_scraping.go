package scraping

import (
	"fmt"
	"net/url"

	"github.com/birdseyeapi/birdseyeapi_v2/internal/models"
)

// SiteScraping handles news scraping from various websites
type SiteScraping struct {
	scrapers      []ScrapingNews
	reactionScrapers []ScrapingReaction
}

// NewSiteScraping creates a new SiteScraping instance
func NewSiteScraping() *SiteScraping {
	// Create a new AI summarizer instance
	summarizer := NewOpenAISummarizer()

	// Create scrapers for each news source
	return &SiteScraping{
		scrapers: []ScrapingNews{
			NewScrapeNewsByCloudWatchImpress(summarizer),
			NewScrapeNewsByHatena(summarizer),
			NewScrapeNewsByZenn(summarizer),
			NewScrapeNewsByZDNet(summarizer),
		},
		reactionScrapers: []ScrapingReaction{
			NewScrapeReactionsByHatena(),
		},
	}
}

// ScrapeNews scrapes news from all configured sources
func (s *SiteScraping) ScrapeNews() ([]models.News, error) {
	var allNews []models.News

	// Scrape news from each source
	for _, scraper := range s.scrapers {
		news, err := scraper.ExtractNews()
		if err != nil {
			fmt.Printf("Error scraping from %s: %v\n", scraper.GetSourceBy(), err)
			continue // Continue with the next source even if one fails
		}
		
		fmt.Printf("%s -> scraped article: %d\n", scraper.GetSourceBy(), len(news))
		allNews = append(allNews, news...)
	}

	return allNews, nil
}

// ScrapeReactions scrapes reactions for a news article
func (s *SiteScraping) ScrapeReactions(news models.News) ([]models.NewsReaction, error) {
	var allReactions []models.NewsReaction

	// Validate the URL before proceeding
	_, err := url.Parse(news.ArticleUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid article URL: %v", err)
	}

	// Scrape reactions from each source
	for _, scraper := range s.reactionScrapers {
		reactions, err := scraper.ExtractReactions(news.ArticleUrl, news.Title)
		if err != nil {
			fmt.Printf("Error scraping reactions from %s: %v\n", scraper.GetSourceBy(), err)
			continue
		}
		
		allReactions = append(allReactions, reactions...)
	}

	return allReactions, nil
}