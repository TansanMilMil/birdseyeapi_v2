package scraping

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

const (
	CloudWatchSourceName      = "CloudWatch by Impress"
	CloudWatchBaseURL         = "https://cloud.watch.impress.co.jp/"
	CloudWatchArticleSelector = "div.article"
	CloudWatchMaxArticles     = 15
)

type ScrapeNewsByCloudWatchImpress struct {
	summarizer Summarizer
}

func NewScrapeNewsByCloudWatchImpress(summarizer Summarizer) *ScrapeNewsByCloudWatchImpress {
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

	doc, err := GetWebDoc(CloudWatchBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := doc.Find(CloudWatchArticleSelector)
	articles = articles.Slice(0, CloudWatchMaxArticles)

	articles.Each(func(i int, selection *goquery.Selection) {
		articleURL := ""
		titleElement := selection.Find("h3 a")
		if href, exists := titleElement.Attr("href"); exists {
			articleURL = href
			if !strings.HasPrefix(articleURL, "http") {
				articleURL = CloudWatchBaseURL + articleURL
			}
		}

		title := strings.TrimSpace(titleElement.Text())

		description := strings.TrimSpace(selection.Find("p.description").Text())

		imageURL := ""
		if src, exists := selection.Find("img").Attr("src"); exists {
			imageURL = src
			if !strings.HasPrefix(imageURL, "http") {
				imageURL = CloudWatchBaseURL + imageURL
			}
		}

		if title != "" && articleURL != "" {
			newsItem := models.News{
				Title:           title,
				Description:     description,
				SourceBy:        CloudWatchSourceName,
				ScrapedUrl:      CloudWatchBaseURL,
				ScrapedDateTime: time.Now(),
				ArticleUrl:      articleURL,
				ArticleImageUrl: imageURL,
			}

			if description != "" && summarizer != nil {
				summary, err := summarizer.Summarize(description)
				if err == nil {
					newsItem.SummarizedText = summary
				}
			}

			news = append(news, newsItem)
			fmt.Print(".")
		}
	})

	return news, nil
}
