package ai

type Summarizer interface {
	Summarize(text string) (string, error)
}
