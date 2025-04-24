package news

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/ai"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/env"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping/doc"
)

const (
	HatenaSourceName      = "Hatena"
	HatenaBaseURL         = "https://b.hatena.ne.jp/hotentry/it"
	HatenaArticleSelector = "#container .entrylist-contents-main"
)

var HatenaMaxArticles = env.GetEnvInt("SCRAPING_ARTICLES", 15)

type ScrapeNewsByHatena struct {
	summarizer ai.Summarizer
}

func NewScrapeNewsByHatena(summarizer ai.Summarizer) *ScrapeNewsByHatena {
	return &ScrapeNewsByHatena{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByHatena) GetSourceBy() string {
	return HatenaSourceName
}

func (s *ScrapeNewsByHatena) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer

	d, err := doc.GetWebDoc(HatenaBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := d.Find(HatenaArticleSelector)
	articles = articles.Slice(0, HatenaMaxArticles)

	articles.Each(func(i int, art *goquery.Selection) {
		titleElement := art.Find(".entrylist-contents-title > a")
		title := titleElement.Text()
		artUrl := titleElement.AttrOr("href", "")

		newsItem := models.News{
			Title:           title,
			Description:     "",
			SourceBy:        HatenaSourceName,
			ScrapedUrl:      HatenaBaseURL,
			ScrapedDateTime: time.Now(),
			ArticleUrl:      artUrl,
			ArticleImageUrl: "",
		}

		artDoc, err := doc.GetWebDoc(artUrl)
		if err != nil {
			fmt.Printf("Failed to parse article HTML: %v\n", err)
			return
		}

		if summarizer != nil {
			summary, err := summarizer.Summarize(artDoc.Text())
			if err == nil {
				newsItem.SummarizedText = summary
			}
		}

		news = append(news, newsItem)
		fmt.Print(".")
	})

	return news, nil
}
