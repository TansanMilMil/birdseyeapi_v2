package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type OpenAISummarizer struct {
	apiKey      string
	baseURL     string
	openAIModel string
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewOpenAISummarizer() *OpenAISummarizer {
	apiKey := os.Getenv("OPENAI_API_KEY")
	baseURL := os.Getenv("OPENAI_CHAT_ENDPOINT")
	openAIModel := os.Getenv("OPENAI_MODEL")

	return &OpenAISummarizer{
		apiKey:      apiKey,
		baseURL:     baseURL,
		openAIModel: openAIModel,
	}
}

func (s *OpenAISummarizer) Summarize(text string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not found")
	}

	reqBody := OpenAIRequest{
		Model: s.openAIModel,
		Messages: []Message{
			{
				Role: "user",
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

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var respBody OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("no summary was generated")
	}

	return respBody.Choices[0].Message.Content, nil
}
