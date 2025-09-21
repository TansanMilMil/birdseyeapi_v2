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
	CloudWatchSourceName      = "CloudWatch by Impress"
	CloudWatchBaseURL         = "https://cloud.watch.impress.co.jp/"
	CloudWatchArticleSelector = "li.item.news"
)

var CloudWatchMaxArticles = env.GetEnvInt("SCRAPING_ARTICLES", 5)

type ScrapeNewsByCloudWatchImpress struct {
	summarizer ai.Summarizer
}

func NewScrapeNewsByCloudWatchImpress(summarizer ai.Summarizer) *ScrapeNewsByCloudWatchImpress {
	return &ScrapeNewsByCloudWatchImpress{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByCloudWatchImpress) GetSourceBy() string {
	return CloudWatchSourceName
}

func (s *ScrapeNewsByCloudWatchImpress) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer

	d, err := doc.GetWebDoc(CloudWatchBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := d.Find(CloudWatchArticleSelector)
	if articles.Length() == 0 {
		return nil, fmt.Errorf("no articles found with selector '%s'", CloudWatchArticleSelector)
	}

	articles = articles.Slice(0, CloudWatchMaxArticles)

	articles.Each(func(i int, art *goquery.Selection) {

		titleElement := art.Find("p.title > a")
		artUrl := strings.TrimSpace(titleElement.AttrOr("href", ""))
		title := strings.TrimSpace(titleElement.Text())

		newsItem := models.News{
			Title:           title,
			Description:     "",
			SourceBy:        CloudWatchSourceName,
			ScrapedUrl:      CloudWatchBaseURL,
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
