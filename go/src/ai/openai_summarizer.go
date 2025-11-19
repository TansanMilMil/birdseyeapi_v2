package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const OPENAI_CHAT_ENDPOINT = "https://api.openai.com/v1/chat/completions"
const OPENAI_MODEL = "gpt-4.1-mini"

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
	apiKey := os.Getenv("BIRDSEYEAPI_V2_OPENAI_API_KEY")
	baseURL := OPENAI_CHAT_ENDPOINT
	openAIModel := OPENAI_MODEL

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
					また、読みやすいように適宜改行を含めてください。
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
		fmt.Print("s.baseURL", s.baseURL)
		fmt.Print("reqBody.Model", reqBody.Model)
		fmt.Print("reqBody.Messages[0].len:", len(reqBody.Messages[0].Content))
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
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
