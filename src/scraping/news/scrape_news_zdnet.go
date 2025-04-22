package news

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/src/ai"
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/src/scraping/doc"
)

const (
	ZDNetSourceName      = "ZDNet Japan"
	ZDNetBaseURL         = "https://japan.zdnet.com/topics/"
	ZDNetArticleSelector = "div.article-item"
	ZDNetMaxArticles     = 15
)

type ScrapeNewsByZDNet struct {
	summarizer ai.Summarizer
}

func NewScrapeNewsByZDNet(summarizer ai.Summarizer) *ScrapeNewsByZDNet {
	return &ScrapeNewsByZDNet{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByZDNet) GetSourceBy() string {
	return ZDNetSourceName
}

func (s *ScrapeNewsByZDNet) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer

	d, err := doc.GetWebDoc(ZDNetBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := d.Find(ZDNetArticleSelector)
	articles = articles.Slice(0, ZDNetMaxArticles)

	articles.Each(func(i int, selection *goquery.Selection) {
		articleURL := ""
		titleElement := selection.Find("h3.title a")
		if href, exists := titleElement.Attr("href"); exists {
			articleURL = href
			if !strings.HasPrefix(articleURL, "http") {
				articleURL = "https://japan.zdnet.com" + articleURL
			}
		}

		title := strings.TrimSpace(titleElement.Text())

		description := strings.TrimSpace(selection.Find("p.abstract").Text())

		imageURL := ""
		if src, exists := selection.Find("div.image img").Attr("src"); exists {
			imageURL = src
		}

		if title != "" && articleURL != "" {
			newsItem := models.News{
				Title:           title,
				Description:     description,
				SourceBy:        ZDNetSourceName,
				ScrapedUrl:      ZDNetBaseURL,
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
