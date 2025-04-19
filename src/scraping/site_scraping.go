package scraping

import (
	"fmt"
	"net/url"

	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

type SiteScraping struct {
	scrapers      []ScrapingNews
	reactionScrapers []ScrapingReaction
}

func NewSiteScraping() *SiteScraping {
	summarizer := NewOpenAISummarizer()

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

func (s *SiteScraping) ScrapeNews() ([]models.News, error) {
	var allNews []models.News

	for _, scraper := range s.scrapers {
		fmt.Print(scraper.GetSourceBy() + ": scraping")
		news, err := scraper.ExtractNews()
		if err != nil {
			fmt.Printf(" -> Error scraping from %s: %v\n", scraper.GetSourceBy(), err)
			continue
		}
		
		fmt.Printf(" -> scraped article: %d\n", scraper.GetSourceBy(), len(news))
		allNews = append(allNews, news...)
	}

	return allNews, nil
}

func (s *SiteScraping) ScrapeReactions(news models.News) ([]models.NewsReaction, error) {
	var allReactions []models.NewsReaction

	_, err := url.Parse(news.ArticleUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid article URL: %v", err)
	}

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