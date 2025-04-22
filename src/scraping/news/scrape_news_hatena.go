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
	HatenaSourceName      = "Hatena"
	HatenaBaseURL         = "https://b.hatena.ne.jp/hotentry/it"
	HatenaArticleSelector = "#container .entrylist-contents-main"
	HatenaMaxArticles     = 15
)

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
		title := strings.TrimSpace(titleElement.Text())
		art_url := HatenaBaseURL + strings.TrimSpace(art.AttrOr("href", ""))

		description := strings.TrimSpace(art.Find("div.entrylist-contents-body").Text())

		newsItem := models.News{
			Title:           title,
			Description:     description,
			SourceBy:        HatenaSourceName,
			ScrapedUrl:      HatenaBaseURL,
			ScrapedDateTime: time.Now(),
			ArticleUrl:      art_url,
			ArticleImageUrl: "",
		}

		art_doc, err := doc.GetWebDoc(art_url)
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
		fmt.Print(".")
	})

	return news, nil
}
