package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("home"))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bucket := "rasc-test-list"

	prefixes := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	numGoroutines := 8
	totalObjects, elapsed, err := fastListS3Objects(context.Background(), s3Client, bucket, prefixes, numGoroutines)

	if err != nil {
		fmt.Printf("Operation completed with errors. First error: %v\n", err)
	}

	fmt.Printf("Total objects found: %d\n", totalObjects)
	fmt.Printf("List operation completed in: %v\n", elapsed)
}

type result struct {
	prefix string
	count  int
	error  error
}

func fastListS3Objects(ctx context.Context, s3Client *s3.Client, bucket string, prefixes []string, numGoroutines int) (int, time.Duration, error) {
	startTime := time.Now()

	maxWorkers := min(len(prefixes), numGoroutines)

	resultChan := make(chan result, len(prefixes))
	workChan := make(chan string, len(prefixes))
	var wg sync.WaitGroup

	for _, prefix := range prefixes {
		workChan <- prefix
	}
	close(workChan)

	for range maxWorkers {
		wg.Go(func() {
			for prefix := range workChan {
				count, err := listObjectsWithPrefix(ctx, s3Client, bucket, prefix)
				resultChan <- result{prefix: prefix, count: count, error: err}
			}
		})
	}

	wg.Wait()
	close(resultChan)

	totalObjects := 0
	var firstError error
	for result := range resultChan {
		if result.error == nil {
			totalObjects += result.count
		} else {
			fmt.Printf("prefix %s: ERROR - %v\n", result.prefix, result.error)
			if firstError == nil {
				firstError = result.error
			}
		}
	}

	elapsed := time.Since(startTime)
	return totalObjects, elapsed, firstError
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
