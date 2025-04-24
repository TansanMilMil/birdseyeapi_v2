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
	ZDNetSourceName      = "ZDNet Japan"
	ZDNetBaseURL         = "https://japan.zdnet.com"
	ZDNetArticleSelector = "#page-wrap > div.pg-container-main > main > section:nth-child(1) > div > ul > li"
)

var ZDNetMaxArticles = env.GetEnvInt("SCRAPING_ARTICLES", 15)

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

	// Check if document is valid
	if d == nil {
		return nil, fmt.Errorf("failed to get document: document is nil")
	}

	articles := d.Find(ZDNetArticleSelector)

	// Check if any articles were found
	if articles.Length() == 0 {
		fmt.Printf("Warning: No articles found with selector '%s'\n", ZDNetArticleSelector)
		return news, nil
	}

	// Slice only if there are enough articles
	if articles.Length() > ZDNetMaxArticles {
		articles = articles.Slice(0, ZDNetMaxArticles)
	}

	articles.Each(func(i int, art *goquery.Selection) {
		// Skip if the article selection is nil
		if art == nil {
			fmt.Println("Warning: Nil article element encountered")
			return
		}

		titleElement := art.Find("a > div.txt > p.txt-ttl")
		if titleElement.Length() == 0 {
			fmt.Println("Warning: Title element not found")
			return
		}

		title := strings.TrimSpace(titleElement.Text())
		if title == "" {
			fmt.Println("Warning: Empty title found")
			return
		}

		artUrlElem := art.Find("a")
		if artUrlElem.Length() == 0 {
			fmt.Println("Warning: Link element not found")
			return
		}

		artUrl := ZDNetBaseURL + strings.TrimSpace(artUrlElem.AttrOr("href", ""))
		if artUrl == ZDNetBaseURL {
			fmt.Println("Warning: Invalid article URL")
			return
		}

		imageURL := ""
		imgElement := art.Find("a > div.thumb > img")
		if imgElement.Length() > 0 {
			if src, exists := imgElement.Attr("src"); exists {
				imageURL = src
			}
		}

		newsItem := models.News{
			Title:           title,
			Description:     "",
			SourceBy:        ZDNetSourceName,
			ScrapedUrl:      ZDNetBaseURL,
			ScrapedDateTime: time.Now(),
			ArticleUrl:      artUrl,
			ArticleImageUrl: imageURL,
		}

		// Only try to get the article doc if we have a valid URL
		if artUrl != ZDNetBaseURL {
			artDoc, err := doc.GetWebDoc(artUrl)
			if err != nil {
				fmt.Printf("Failed to parse article HTML: %v\n", err)
				return
			}

			if artDoc != nil && summarizer != nil {
				summary, err := summarizer.Summarize(artDoc.Text())
				if err == nil {
					newsItem.SummarizedText = summary
				} else {
					fmt.Printf("Failed to summarize article: %v\n", err)
				}
			}
		}

		news = append(news, newsItem)
		fmt.Print(".")
	})

	return news, nil
}
