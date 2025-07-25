package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type result struct {
	prefix string
	count  int
	error  error
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("home"))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bucket := "rasc-test-list"

	prefixes := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

	startTime := time.Now()

	resultChan := make(chan result, len(prefixes))
	var wg sync.WaitGroup

	for _, prefix := range prefixes {
		wg.Add(1)
		go func(prefix string) {
			defer wg.Done()
			count, err := listObjectsWithPrefix(context.Background(), s3Client, bucket, prefix)
			resultChan <- result{prefix: prefix, count: count, error: err}
		}(prefix)
	}

	wg.Wait()
	close(resultChan)

	totalObjects := 0
	for result := range resultChan {
		if result.error == nil {
			totalObjects += result.count
		} else {
			fmt.Printf("prefix %s: ERROR - %v\n", result.prefix, result.error)
		}
	}

	elapsed := time.Since(startTime)

	fmt.Printf("Total objects found: %d\n", totalObjects)
	fmt.Printf("List operation completed in: %v\n", elapsed)
}

func listObjectsWithPrefix(ctx context.Context, s3Client *s3.Client, bucket, prefix string) (int, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	paginator := s3.NewListObjectsV2Paginator(s3Client, input)
	var objectCount int

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return 0, err
		}

		for range page.Contents {
			objectCount++
		}
	}

	return objectCount, nil
}
