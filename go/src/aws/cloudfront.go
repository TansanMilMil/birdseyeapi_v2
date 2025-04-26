package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/birdseyeapi/birdseyeapi_v2/go/src/env"
)

var (
	AWSRegion = env.GetEnv("AWS_REGION", "")
)

// CreateInvalidation creates a CloudFront invalidation for the specified paths
func CreateInvalidation(distributionID string, paths []string) error {
	if distributionID == "" {
		return fmt.Errorf("AWS_CLOUDFRONT_BIRDSEYEAPIPROXY_DISTRIBUTION_ID environment variable not set")
	}
	// Create AWS session
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWSRegion),
	)

	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create CloudFront client
	svc := cloudfront.NewFromConfig(cfg)

	// Prepare invalidation input
	callerReference := fmt.Sprintf("birdseyeapi-invalidation-%d", time.Now().Unix())

	items := make([]string, len(paths))
	copy(items, paths)

	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionID),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(callerReference),
			Paths: &types.Paths{
				Quantity: aws.Int32(int32(len(items))),
				Items:    items,
			},
		},
	}

	_, err = svc.CreateInvalidation(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to create CloudFront invalidation: %v", err)
	}

	fmt.Printf("CloudFront invalidation created successfully with reference: %s\n", callerReference)
	return nil
}
