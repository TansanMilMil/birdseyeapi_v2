package scraping

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/src/models"
)

type ScrapeNewsByZenn struct {
	summarizer Summarizer
}

func NewScrapeNewsByZenn(summarizer Summarizer) *ScrapeNewsByZenn {
	return &ScrapeNewsByZenn{
		summarizer: summarizer,
	}
}

func (s *ScrapeNewsByZenn) GetSourceBy() string {
	return "Zenn"
}

func (s *ScrapeNewsByZenn) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer

	url := "https://zenn.dev"
	doc, err := GetWebDoc(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	articles := doc.Find("#tech-trend > div > div > div > article > div > a[class^=\"ArticleList_link\"]")
	articles.Each(func(i int, art *goquery.Selection) {
		title := strings.TrimSpace(art.Find("h2").Text())
		art_url := url + strings.TrimSpace(art.AttrOr("href", ""))
          
		newsItem := models.News{
			Title:           title,
			Description:     "",
			SourceBy:        s.GetSourceBy(),
			ScrapedUrl:      url,
			ScrapedDateTime: time.Now(),
			ArticleUrl:      art_url,
			ArticleImageUrl: "",
		}

		art_doc, err := GetWebDoc(art_url)
		if err != nil {
			fmt.Printf("Failed to parse article HTML: %v\n", err)
			return
		}

		if summarizer != nil {
			summary, err := summarizer.Summarize(art_doc.Text())
			if err == nil {
				newsItem.SummarizedText = summary
			}
		}

		news = append(news, newsItem)
	})

	return news, nil
}