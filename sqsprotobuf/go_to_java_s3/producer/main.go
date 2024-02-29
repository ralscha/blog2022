package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"log"
	"producer/shared"
)

const queueURL = "https://sqs.us-east-1.amazonaws.com/660461151343/queue-d8494a4"
const messageBucket = "messages-a5b0326"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("home"))
	check(err)

	sqsClient := sqs.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)

	var people []*shared.Person
	for i := 1; i < 10_000; i++ {
		p := &shared.Person{
			Id:    int32(i),
			Name:  faker.Name(),
			Email: faker.Email(),
			Phones: []*shared.Person_PhoneNumber{
				{Number: faker.Phonenumber(), Type: shared.Person_MOBILE},
			}}

		people = append(people, p)
	}
	book := &shared.AddressBook{People: people}

	out, err := proto.Marshal(book)
	check(err)

	s3Key, err := uuid.NewUUID()
	check(err)

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(messageBucket),
		Key:    aws.String(s3Key.String()),
		Body:   bytes.NewReader(out),
	})
	check(err)

	msg := &sqs.SendMessageInput{
		MessageBody: aws.String(s3Key.String()),
		QueueUrl:    aws.String(queueURL),
	}
	message, err := sqsClient.SendMessage(context.TODO(), msg)
	check(err)
	fmt.Println(message.MessageId)
}

func check(e error) {
	if e != nil {
		log.Panicln(e)
	}
}
