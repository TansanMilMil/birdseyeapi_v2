package scraping

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

const (
	HatenaSourceName      = "Hatena"
	HatenaBaseURL         = "https://b.hatena.ne.jp/hotentry/it"
	HatenaArticleSelector = "div.entrylist-item"
	HatenaMaxArticles     = 15
)

type ScrapeNewsByHatena struct {
	summarizer Summarizer
}

func NewScrapeNewsByHatena(summarizer Summarizer) *ScrapeNewsByHatena {
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

	doc, err := GetWebDoc(HatenaBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := doc.Find(HatenaArticleSelector)
	articles = articles.Slice(0, HatenaMaxArticles)

	articles.Each(func(i int, selection *goquery.Selection) {
		articleURL := ""
		titleElement := selection.Find("h3.entrylist-contents-title a")
		if href, exists := titleElement.Attr("href"); exists {
			articleURL = href
		}

		title := strings.TrimSpace(titleElement.Text())

		description := strings.TrimSpace(selection.Find("div.entrylist-contents-body").Text())

		imageURL := ""
		if src, exists := selection.Find("img.entrylist-contents-thumb").Attr("src"); exists {
			imageURL = src
		}

		if title != "" && articleURL != "" {
			newsItem := models.News{
				Title:           title,
				Description:     description,
				SourceBy:        HatenaSourceName,
				ScrapedUrl:      HatenaBaseURL,
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
