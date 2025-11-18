package ai

import (
	"os"
	"testing"
)

func TestNewClaudeSummarizer(t *testing.T) {
	// 環境変数を保存
	originalAPIKey := os.Getenv("BIRDSEYEAPI_V2_CLAUDE_API_KEY")
	defer func() {
		os.Setenv("BIRDSEYEAPI_V2_CLAUDE_API_KEY", originalAPIKey)
	}()

	t.Run("環境変数が設定されている場合", func(t *testing.T) {
		os.Setenv("BIRDSEYEAPI_V2_CLAUDE_API_KEY", "test-api-key")

		summarizer := NewClaudeSummarizer()

		if summarizer.apiKey != "test-api-key" {
			t.Errorf("Expected apiKey to be 'test-api-key', got '%s'", summarizer.apiKey)
		}
		if summarizer.maxTokens != 1024 {
			t.Errorf("Expected maxTokens to be 1024, got %d", summarizer.maxTokens)
		}
	})

	t.Run("環境変数が設定されていない場合", func(t *testing.T) {
		os.Unsetenv("BIRDSEYEAPI_V2_CLAUDE_API_KEY")

		summarizer := NewClaudeSummarizer()

		if summarizer.apiKey != "" {
			t.Errorf("Expected apiKey to be empty, got '%s'", summarizer.apiKey)
		}
		if summarizer.baseURL != CLAUDE_CHAT_ENDPOINT {
			t.Errorf("Expected baseURL to be '%s', got '%s'", CLAUDE_CHAT_ENDPOINT, summarizer.baseURL)
		}
		if summarizer.claudeModel != CLAUDE_MODEL {
			t.Errorf("Expected claudeModel to be '%s', got '%s'", CLAUDE_MODEL, summarizer.claudeModel)
		}
	})
}

func TestClaudeSummarizer_Summarize_NoAPIKey(t *testing.T) {
	summarizer := &ClaudeSummarizer{
		apiKey:      "",
		baseURL:     "https://api.anthropic.com/v1/messages",
		claudeModel: "claude-3-5-sonnet-20241022",
		maxTokens:   1024,
	}

	_, err := summarizer.Summarize("test text")
	if err == nil {
		t.Error("Expected error when API key is not set")
	}
	if err.Error() != "claude API key not found" {
		t.Errorf("Expected error message 'claude API key not found', got '%s'", err.Error())
	}
}

// インテグレーションテスト: 実際のClaude APIを呼び出す
// 環境変数 BIRDSEYEAPI_V2_CLAUDE_API_KEY が設定されている場合のみ実行される
func TestClaudeSummarizer_Integration(t *testing.T) {
	t.Skip("Skipping integration test: not ready for execution") // 一時的にスキップ

	apiKey := os.Getenv("BIRDSEYEAPI_V2_CLAUDE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test: BIRDSEYEAPI_V2_CLAUDE_API_KEY not set")
	}

	summarizer := NewClaudeSummarizer()

	testText := `Go言語は、Googleによって開発されたプログラミング言語です。
シンプルで効率的なコードを書くことができ、並行処理に優れています。
静的型付けとガベージコレクションを備えており、
モダンなソフトウェア開発に適した言語として広く使用されています。
特に、サーバーサイドのアプリケーション、クラウドインフラストラクチャ、
コマンドラインツールなどの開発に人気があります。`

	t.Log("Calling Claude API for summarization...")
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
