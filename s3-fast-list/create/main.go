package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
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

	const numGoroutines = 20
	const totalObjects = 25_000
	objectsPerGoroutine := totalObjects / numGoroutines

	var wg sync.WaitGroup
	errChan := make(chan error, totalObjects)

	fmt.Printf("Starting to create %d objects with %d goroutines...\n", totalObjects, numGoroutines)

	for range numGoroutines {
		wg.Go(func() {
			for range objectsPerGoroutine {
				key := uuid.New().String()
				_, err := s3Client.PutObject(context.Background(), &s3.PutObjectInput{
					Bucket: &bucket,
					Key:    &key,
					Body:   nil, // empty content
				})
				if err != nil {
					errChan <- err
				}
			}
		})
	}

	wg.Wait()
	close(errChan)

	var errCount int
	for err := range errChan {
		if err != nil {
			log.Printf("error creating object: %v", err)
			errCount++
		}
	}

	if errCount > 0 {
		fmt.Printf("Finished with %d errors.\n", errCount)
	} else {
		fmt.Println("Finished creating all objects successfully.")
	}
}
