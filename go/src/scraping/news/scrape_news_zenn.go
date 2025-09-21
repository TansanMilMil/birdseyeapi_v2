package news

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/ai"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/env"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/scraping/doc"
)

const (
	ZennSourceName      = "Zenn"
	ZennBaseURL         = "https://zenn.dev"
	ZennArticleSelector = "#tech-trend > div > div > div > article > div > a"
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
	return ZennSourceName
}

func (s *ScrapeNewsByZenn) ExtractNews() ([]models.News, error) {
	var news []models.News
	summarizer := s.summarizer

	d, err := doc.GetWebDoc(ZennBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Zennは <script id="__NEXT_DATA__">...</script>の"..."にJSON文字列で記事データが埋め込まれているのでうまく取り出す
	var articles []map[string]interface{}

	d.Find("script#__NEXT_DATA__").First().Each(func(i int, s *goquery.Selection) {
		scriptContent := s.Text()

		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(scriptContent), &jsonData); err != nil {
			fmt.Printf("Failed to parse JSON: %v\n", err)
			return
		}

		// props.pageProps.dailyTechArticlesにアクセス
		if props, ok := jsonData["props"].(map[string]interface{}); ok {
			if pageProps, ok := props["pageProps"].(map[string]interface{}); ok {
				if dailyTechArticles, ok := pageProps["dailyTechArticles"].([]interface{}); ok {
					for _, article := range dailyTechArticles {
						if articleMap, ok := article.(map[string]interface{}); ok {
							articles = append(articles, articleMap)
						}
					}
				}
			}
		}
	})

	if len(articles) == 0 {
		return nil, fmt.Errorf("no articles found in __NEXT_DATA__")
	}

	// MaxArticlesの制限を適用
	if len(articles) > MaxArticles {
		articles = articles[:MaxArticles]
	}

	for _, article := range articles {
		title, titleOk := article["title"].(string)
		path, pathOk := article["path"].(string)

		if !titleOk || !pathOk {
			continue
		}

		artUrl := ZennBaseURL + path

		newsItem := models.News{
			Title:           title,
			Description:     "",
			SourceBy:        s.GetSourceBy(),
			ScrapedUrl:      ZennBaseURL,
			ScrapedDateTime: time.Now(),
			ArticleUrl:      artUrl,
			ArticleImageUrl: "",
		}

		art_doc, err := doc.GetWebDoc(artUrl)
		if err != nil {
			fmt.Printf("Failed to parse article HTML: %v\n", err)
			continue
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
	}

	return news, nil
}
