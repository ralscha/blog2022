package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"google.golang.org/protobuf/proto"
	"log"
	"producer/shared"
)

const queueURL = "https://sqs.us-east-1.amazonaws.com/660461151343/queue-d8494a4"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("home"))
	check(err)

	sqsClient := sqs.NewFromConfig(cfg)

	p := shared.Person{
		Id:    1,
		Name:  "John",
		Email: "john@test.com",
		Phones: []*shared.Person_PhoneNumber{
			{Number: "1111", Type: shared.Person_MOBILE},
		},
	}
	book := &shared.AddressBook{People: []*shared.Person{&p}}
	out, err := proto.Marshal(book)
	check(err)

	msg := &sqs.SendMessageInput{
		MessageBody: aws.String(" "),
		QueueUrl:    aws.String(queueURL),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Body": {
				DataType:    aws.String("Binary"),
				BinaryValue: out,
			},
		},
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
