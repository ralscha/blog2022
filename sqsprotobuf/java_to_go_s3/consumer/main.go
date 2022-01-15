package main

import (
	"bytes"
	"consumer/shared"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

const queueURL = "https://sqs.us-east-1.amazonaws.com/660461151343/queue-d8494a4"
const messageBucket = "messages-a5b0326"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("home"))
	check(err)

	sqsClient := sqs.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)

	for {
		fmt.Println(time.Now(), ": poll for message")
		receiveMessageInput := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20,
		}
		receiveMessageOutput, err := sqsClient.ReceiveMessage(context.TODO(), receiveMessageInput)
		check(err)

		for _, message := range receiveMessageOutput.Messages {
			s3Key := message.Body

			getObjectInput := &s3.GetObjectInput{
				Bucket: aws.String(messageBucket),
				Key:    s3Key,
			}
			object, err := s3Client.GetObject(context.TODO(), getObjectInput)
			check(err)

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(object.Body)
			check(err)

			ab := &shared.AddressBook{}
			err = proto.Unmarshal(buf.Bytes(), ab)
			check(err)

			for _, person := range ab.People {
				fmt.Println(person.Id)
				fmt.Println(person.Name)
				fmt.Println(person.Email)
				for _, phone := range person.Phones {
					fmt.Print(phone.Type)
					fmt.Print(" : ")
					fmt.Println(phone.Number)
				}
				fmt.Println()
			}

			deleteObjectInput := &s3.DeleteObjectInput{
				Bucket: aws.String(messageBucket),
				Key:    s3Key,
			}
			_, err = s3Client.DeleteObject(context.TODO(), deleteObjectInput)
			check(err)

			deleteMessageInput := &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: message.ReceiptHandle,
			}
			_, err = sqsClient.DeleteMessage(context.TODO(), deleteMessageInput)
			check(err)
		}
	}

}

func check(e error) {
	if e != nil {
		log.Panicln(e)
	}
}
