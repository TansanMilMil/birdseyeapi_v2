package scraping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// OpenAISummarizer implements the Summarizer interface using OpenAI API
type OpenAISummarizer struct {
	apiKey  string
	baseURL string
	openAIModel string
}

// OpenAIRequest represents a request to the OpenAI API
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message represents a message in an OpenAI chat conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents a response from the OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewOpenAISummarizer creates a new OpenAISummarizer
func NewOpenAISummarizer() *OpenAISummarizer {
	apiKey := os.Getenv("OPENAI_API_KEY")
	baseURL := os.Getenv("OPENAI_CHAT_ENDPOINT")
	openAIModel := os.Getenv("OPENAI_MODEL")
	
	return &OpenAISummarizer{
		apiKey:  apiKey,
		baseURL: baseURL,
		openAIModel: openAIModel,
	}
}

// Summarize summarizes the given text using OpenAI API
func (s *OpenAISummarizer) Summarize(text string) (string, error) {
	// Check if API key is available
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not found")
	}

	// Prepare request body
	reqBody := OpenAIRequest{
		Model: s.openAIModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf(`次の文章を日本語で要約してください。
                    なお、要約結果の文章は200文字以内に収まるように調整してください。
                    ---
                    %s`, text),
			},
		},
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	
	// Add headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.apiKey)
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status code %d", resp.StatusCode)
	}
	
	// Parse response
	var respBody OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}
	
	// Extract summarized text
	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("no summary was generated")
	}
	
	return respBody.Choices[0].Message.Content, nil
}