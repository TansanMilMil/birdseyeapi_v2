package cache

type CDNInvalidator interface {
	Invalidate() bool
}
