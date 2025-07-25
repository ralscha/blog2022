package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// on AWS
	// cfg, err := config.LoadDefaultConfig(context.TODO())

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("home"))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bucket := "rasc-test-list"

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	paginator := s3.NewListObjectsV2Paginator(s3Client, input)

	startTime := time.Now()

	var objectCount int
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			log.Fatalf("failed to get page: %v", err)
		}

		for range page.Contents {
			objectCount++
		}
	}

	elapsed := time.Since(startTime)

	fmt.Printf("\nTotal objects found: %d\n", objectCount)
	fmt.Printf("List operation completed in: %v\n", elapsed)
}
