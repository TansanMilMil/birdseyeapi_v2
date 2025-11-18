package ai

import (
	"os"
	"testing"
)

func TestNewOpenAISummarizer(t *testing.T) {
	// 環境変数を保存
	originalAPIKey := os.Getenv("BIRDSEYEAPI_V2_OPENAI_API_KEY")
	defer func() {
		os.Setenv("BIRDSEYEAPI_V2_OPENAI_API_KEY", originalAPIKey)
	}()

	t.Run("環境変数が設定されている場合", func(t *testing.T) {
		os.Setenv("BIRDSEYEAPI_V2_OPENAI_API_KEY", "test-api-key")

		summarizer := NewOpenAISummarizer()

		if summarizer.apiKey != "test-api-key" {
			t.Errorf("Expected apiKey to be 'test-api-key', got '%s'", summarizer.apiKey)
		}
		if summarizer.baseURL != OPENAI_CHAT_ENDPOINT {
			t.Errorf("Expected baseURL to be '%s', got '%s'", OPENAI_CHAT_ENDPOINT, summarizer.baseURL)
		}
		if summarizer.openAIModel != OPENAI_MODEL {
			t.Errorf("Expected openAIModel to be '%s', got '%s'", OPENAI_MODEL, summarizer.openAIModel)
		}
	})

	t.Run("環境変数が設定されていない場合", func(t *testing.T) {
		os.Unsetenv("BIRDSEYEAPI_V2_OPENAI_API_KEY")

		summarizer := NewOpenAISummarizer()

		if summarizer.apiKey != "" {
			t.Errorf("Expected apiKey to be empty, got '%s'", summarizer.apiKey)
		}
		if summarizer.baseURL != OPENAI_CHAT_ENDPOINT {
			t.Errorf("Expected baseURL to be '%s', got '%s'", OPENAI_CHAT_ENDPOINT, summarizer.baseURL)
		}
		if summarizer.openAIModel != OPENAI_MODEL {
			t.Errorf("Expected openAIModel to be '%s', got '%s'", OPENAI_MODEL, summarizer.openAIModel)
		}
	})
}

func TestOpenAISummarizer_Summarize_NoAPIKey(t *testing.T) {
	summarizer := &OpenAISummarizer{
		apiKey:      "",
		baseURL:     "https://api.openai.com/v1/chat/completions",
		openAIModel: "gpt-4.1-mini",
	}

	_, err := summarizer.Summarize("test text")
	if err == nil {
		t.Error("Expected error when API key is not set")
	}
	if err.Error() != "OpenAI API key not found" {
		t.Errorf("Expected error message 'OpenAI API key not found', got '%s'", err.Error())
	}
}

// インテグレーションテスト: 実際のOpenAI APIを呼び出す
// 環境変数 BIRDSEYEAPI_V2_OPENAI_API_KEY が設定されている場合のみ実行される
func TestOpenAISummarizer_Integration(t *testing.T) {
	apiKey := os.Getenv("BIRDSEYEAPI_V2_OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: BIRDSEYEAPI_V2_OPENAI_API_KEY not set")
	}

	summarizer := NewOpenAISummarizer()

	testText := `Go言語は、Googleによって開発されたプログラミング言語です。
シンプルで効率的なコードを書くことができ、並行処理に優れています。
静的型付けとガベージコレクションを備えており、
モダンなソフトウェア開発に適した言語として広く使用されています。
特に、サーバーサイドのアプリケーション、クラウドインフラストラクチャ、
コマンドラインツールなどの開発に人気があります。`

	t.Log("Calling OpenAI API for summarization...")
	summary, err := summarizer.Summarize(testText)
	if err != nil {
		t.Fatalf("Failed to summarize: %v", err)
	}

	t.Logf("Summary: %s", summary)

	// 要約が空でないことを確認
	if summary == "" {
		t.Error("Expected non-empty summary")
	}
}
