package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const CLAUDE_CHAT_ENDPOINT = "https://api.anthropic.com/v1/messages"
const CLAUDE_MODEL = "claude-3-5-sonnet-20241022"

type ClaudeSummarizer struct {
	apiKey      string
	baseURL     string
	claudeModel string
	maxTokens   int
}

type ClaudeRequest struct {
	Model     string          `json:"model"`
	Messages  []ClaudeMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func NewClaudeSummarizer() *ClaudeSummarizer {
	apiKey := os.Getenv("BIRDSEYEAPI_V2_CLAUDE_API_KEY")
	baseURL := CLAUDE_CHAT_ENDPOINT
	claudeModel := CLAUDE_MODEL

	return &ClaudeSummarizer{
		apiKey:      apiKey,
		baseURL:     baseURL,
		claudeModel: claudeModel,
		maxTokens:   1024,
	}
}

func (s *ClaudeSummarizer) Summarize(text string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("claude API key not found")
	}

	reqBody := ClaudeRequest{
		Model: s.claudeModel,
		Messages: []ClaudeMessage{
			{
				Role: "user",
				Content: fmt.Sprintf(`次の文章を日本語で要約してください。
                    なお、要約結果の文章は200文字以内に収まるように調整してください。
					また、読みやすいように適宜改行を含めてください。
                    ---
                    %s`, text),
			},
		},
		MaxTokens: s.maxTokens,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", s.apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var respBody ClaudeResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(respBody.Content) == 0 {
		return "", fmt.Errorf("no summary was generated")
	}

	return respBody.Content[0].Text, nil
}
