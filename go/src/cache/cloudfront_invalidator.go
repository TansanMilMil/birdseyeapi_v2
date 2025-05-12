package cache

import (
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/aws"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/env"
)

type CloudFrontInvalidator struct{}

func (c *CloudFrontInvalidator) Invalidate() bool {
	err := aws.CreateInvalidation(
		env.GetEnv("AWS_CLOUDFRONT_BIRDSEYEAPIPROXY_DISTRIBUTION_ID", ""),
		[]string{
			"/news/today-news",
			"/news/news-reactions/*",
			"/news/trends",
		})
	if err != nil {
		println("Error creating CloudFront invalidation:", err.Error())
	} else {
		println("CloudFront invalidation created successfully")
	}

	return err == nil
}
