package api

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/birdseyeapi/birdseyeapi_v2/go/src/models"
	"github.com/gin-gonic/gin"
)

type MockNewsRepository struct {
	news []models.News
}

func (m *MockNewsRepository) GetNews(t time.Time, c *gin.Context) []models.News {
	var result []models.News
	targetDate := t.Format("2006-01-02")
	for _, n := range m.news {
		newsDate := n.CreatedAt.Format("2006-01-02")
		if newsDate >= targetDate {
			result = append(result, n)
		}
	}
	return result
}

func setupTestHandler(mockNews []models.News) *NewsHandler {
	handler := &NewsHandler{
		newsRepo: &MockNewsRepository{news: mockNews},
	}
	return handler
}

func createMockNews(date time.Time, count int) []models.News {
	var news []models.News
	for i := 0; i < count; i++ {
		n := models.News{
			Title:           "Test News",
			Description:     "Test Description",
			SourceBy:        "Test Source",
			ScrapedUrl:      "https://example.com",
			ScrapedDateTime: date,
			ArticleUrl:      "https://example.com/article",
		}
		n.CreatedAt = date
		news = append(news, n)
	}
	return news
}

func TestSearchNewsWithBackoff_NewsFoundOnFirstDay(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	today := time.Now()
	mockNews := createMockNews(today, 5)
	handler := setupTestHandler(mockNews)

	result, err := handler.SearchNewsWithBackoff(10, c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected news to be found, but got empty result")
	}

	if len(result) != 5 {
		t.Errorf("Expected 5 news items, got %d", len(result))
	}
}

func TestSearchNewsWithBackoff_NoNewsFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var mockNews []models.News
	handler := setupTestHandler(mockNews)

	result, err := handler.SearchNewsWithBackoff(5, c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected no news to be found, but got %d items", len(result))
	}
}

func TestSearchNewsWithBackoff_StopsWhenNewsFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	twoDaysAgo := time.Now().AddDate(0, 0, -2)
	fiveDaysAgo := time.Now().AddDate(0, 0, -5)

	var mockNews []models.News
	mockNews = append(mockNews, createMockNews(twoDaysAgo, 1)...)
	mockNews = append(mockNews, createMockNews(fiveDaysAgo, 3)...)
	handler := setupTestHandler(mockNews)

	result, err := handler.SearchNewsWithBackoff(10, c)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected function to stop at first news found (1 item), but got %d items", len(result))
	}
}

func TestSearchNewsWithBackoff_NegativeBackoffDays(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	today := time.Now()
	mockNews := createMockNews(today, 1)
	handler := setupTestHandler(mockNews)

	_, err := handler.SearchNewsWithBackoff(-1, c)

	if err == nil {
		t.Errorf("Expected error for negative backoffDays, but got none")
	}
}
