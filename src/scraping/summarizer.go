package scraping

type Summarizer interface {
	Summarize(text string) (string, error)
}