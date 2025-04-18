package scraping

// Summarizer defines the interface for text summarization
type Summarizer interface {
	// Summarize summarizes the given text
	Summarize(text string) (string, error)
}