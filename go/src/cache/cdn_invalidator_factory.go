package cache

type CDNInvalidatorFactory struct{}

func (f *CDNInvalidatorFactory) CreateInvalidator() CDNInvalidator {
	return &CloudFrontInvalidator{}
}
