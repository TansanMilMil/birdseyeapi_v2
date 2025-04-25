package scraping

import (
	"fmt"
	"net/url"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/ai"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping/news"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping/reaction"
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
	var allNews []models.News

	for _, scraper := range s.scrapers {
		fmt.Print(scraper.GetSourceBy() + ": scraping")
		news, err := scraper.ExtractNews()
		if err != nil {
			fmt.Printf(" -> Error scraping from %s: %v\n", scraper.GetSourceBy(), err)
			continue
		}

		fmt.Printf(" -> scraped article: %s: %d\n", scraper.GetSourceBy(), len(news))
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
		reactions, err := scraper.ExtractReactions(news.ID, news.ArticleUrl, news.Title)
		if err != nil {
			fmt.Printf("Error scraping reactions from %s: %v\n", scraper.GetSourceBy(), err)
			continue
		}

		fmt.Printf(" -> scraped reactions: %s: %d\n", scraper.GetSourceBy(), len(reactions))
		allReactions = append(allReactions, reactions...)
	}

	return allReactions, nil
}
