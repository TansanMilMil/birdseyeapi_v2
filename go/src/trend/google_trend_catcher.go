package trend

import (
	"fmt"
	"time"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/mmcdole/gofeed"
)

const GoogleRSSTrendsDaily string = "https://trends.google.co.jp/trending/rss?geo=JP"

type GoogleTrendCatcher struct {
}

func (gt *GoogleTrendCatcher) GetTrends() ([]models.News, error) {
	var newsList []models.News
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(GoogleRSSTrendsDaily)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	for _, entry := range feed.Items {
		news := entryToNews(entry)
		newsList = append(newsList, news)
	}

	return newsList, nil
}

func entryToNews(entry *gofeed.Item) models.News {
	news := models.News{
		Title:           entry.Title,
		ScrapedUrl:      entry.Link,
		SourceBy:        "googleTrends",
		ScrapedDateTime: time.Now().UTC(),
	}

	if extensions, ok := entry.Extensions["ht"]; ok {
		if newsItems, ok := extensions["news_item"]; ok {
			for _, newsItem := range newsItems {
				if newsItem.Children != nil {
					if snippets, ok := newsItem.Children["news_item_snippet"]; ok && len(snippets) > 0 {
						news.Description = snippets[0].Value
					}

					if urls, ok := newsItem.Children["news_item_url"]; ok && len(urls) > 0 {
						news.ArticleUrl = urls[0].Value
					}
				}
			}
		}

		if pictures, ok := extensions["picture"]; ok && len(pictures) > 0 {
			news.ArticleImageUrl = pictures[0].Value
		}
	}
	return news
}
