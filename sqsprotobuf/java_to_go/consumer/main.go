package main

import (
	"consumer/shared"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

const queueURL = "https://sqs.us-east-1.amazonaws.com/660461151343/queue-d8494a4"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("home"))
	check(err)

	sqsClient := sqs.NewFromConfig(cfg)

	for {
		fmt.Println(time.Now(), ": poll for message")
		receiveMessageInput := &sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(queueURL),
			MaxNumberOfMessages:   10,
			WaitTimeSeconds:       20,
			MessageAttributeNames: []string{"Body"},
		}
		receiveMessageOutput, err := sqsClient.ReceiveMessage(context.TODO(), receiveMessageInput)
		check(err)

		for _, message := range receiveMessageOutput.Messages {
			body := message.MessageAttributes["Body"]

			ab := &shared.AddressBook{}
			err := proto.Unmarshal(body.BinaryValue, ab)
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
