package models

import (
	"time"
)

// NewsResponse represents the API response format for news articles
type GetAllNewsResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	SummarizedText  string    `json:"summarizedText"`
	SourceBy        string    `json:"sourceBy"`
	ScrapedUrl      string    `json:"scrapedUrl"`
	ScrapedDateTime time.Time `json:"scrapedDateTime"`
	ArticleUrl      string    `json:"articleUrl"`
	ArticleImageUrl string    `json:"articleImageUrl"`
	ReactionCount   uint      `json:"reactionCount"`
}

func ToGetAllNewsResponse(n []News) []GetAllNewsResponse {
	results := make([]GetAllNewsResponse, len(n))
	for i, news := range n {
		results[i] = GetAllNewsResponse{
			ID:              news.ID,
			Title:           news.Title,
			Description:     news.Description,
			SummarizedText:  news.SummarizedText,
			SourceBy:        news.SourceBy,
			ScrapedUrl:      news.ScrapedUrl,
			ScrapedDateTime: news.ScrapedDateTime,
			ArticleUrl:      news.ArticleUrl,
			ArticleImageUrl: news.ArticleImageUrl,
			ReactionCount:   uint(len(news.Reactions)),
		}
	}

	return results
}
