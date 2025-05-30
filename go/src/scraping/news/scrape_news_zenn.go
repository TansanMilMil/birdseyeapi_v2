package news

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/ai"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/env"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping/doc"
)

const (
	SourceName      = "Zenn"
	BaseURL         = "https://zenn.dev"
	ArticleSelector = "#tech-trend > div > div > div > article > div > a[class^=\"ArticleList_link\"]"
)

var MaxArticles = env.GetEnvInt("SCRAPING_ARTICLES", 15)

type ScrapeNewsByZenn struct {
	summarizer ai.Summarizer
}

func NewScrapeNewsByZenn(summarizer ai.Summarizer) *ScrapeNewsByZenn {
	return &ScrapeNewsByZenn{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByZenn) GetSourceBy() string {
	return SourceName
}

func (s *ScrapeNewsByZenn) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer

	d, err := doc.GetWebDoc(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := d.Find(ArticleSelector)
	articles = articles.Slice(0, MaxArticles)

	articles.Each(func(i int, art *goquery.Selection) {
		title := strings.TrimSpace(art.Find("h2").Text())
		artUrl := BaseURL + strings.TrimSpace(art.AttrOr("href", ""))

		newsItem := models.News{
			Title:           title,
			Description:     "",
			SourceBy:        s.GetSourceBy(),
			ScrapedUrl:      BaseURL,
			ScrapedDateTime: time.Now(),
			ArticleUrl:      artUrl,
			ArticleImageUrl: "",
		}

		art_doc, err := doc.GetWebDoc(artUrl)
		if err != nil {
			fmt.Printf("Failed to parse article HTML: %v\n", err)
			return
		}

		if summarizer != nil {
			summary, err := summarizer.Summarize(art_doc.Text())
			if err != nil {
				fmt.Printf("Failed to summarize article: %v\n", err)
			} else {
				newsItem.SummarizedText = summary
			}
		}

		news = append(news, newsItem)
		fmt.Print(".")
	})

	return news, nil
}
