package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"os"
)

var dynamodbClient *dynamodb.Client
var tableName = os.Getenv("TABLE_NAME")
var errorResponse = events.APIGatewayV2HTTPResponse{
	StatusCode: 500,
	Body:       "Internal Error",
}

type Todo struct {
	Id          string `json:"id"`
	DueDate     string `json:"dueDate,omitempty"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type TodoPostResponse struct {
	FieldErrors map[string]string `json:"fieldErrors"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("init: unable to load AWS config, %v", err)
	}
	dynamodbClient = dynamodb.NewFromConfig(cfg)
}

func main() {
	lambda.Start(handle)
}

func handle(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch request.RouteKey {
	case "GET /todos":
		return getTodos()
	case "POST /todos":
		return updateTodo(request)
	case "DELETE /todos/{id}":
		return deleteTodo(request)
	}

	return errorResponse, nil
}

func getTodos() (events.APIGatewayV2HTTPResponse, error) {
	p := dynamodb.NewScanPaginator(dynamodbClient, &dynamodb.ScanInput{
		TableName: &tableName,
	})

	todos := make([]Todo, 0)
	for p.HasMorePages() {
		out, err := p.NextPage(context.Background())
		if err != nil {
			return errorResponse, err
		}

		var pTodos []Todo
		err = attributevalue.UnmarshalListOfMaps(out.Items, &pTodos)
		if err != nil {
			return errorResponse, err
		}
		todos = append(todos, pTodos...)
	}

	j, err := json.Marshal(todos)
	if err != nil {
		return errorResponse, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"content-type": "application/json"},
		Body:       string(j),
	}, nil

}

func updateTodo(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var todo Todo
	err := json.Unmarshal([]byte(request.Body), &todo)
	if err != nil {
		return errorResponse, err
	}

	errors := make(map[string]string)
	if todo.Id == "" {
		errors["id"] = "required"
	}
	if todo.Description == "" {
		errors["description"] = "required"
	}
	if todo.Priority == "" {
		errors["priority"] = "required"
	}
	if len(errors) > 0 {
		j, err := json.Marshal(TodoPostResponse{
			FieldErrors: errors,
		})
		if err != nil {
			return errorResponse, err
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: 422,
			Headers:    map[string]string{"content-type": "application/json"},
			Body:       string(j),
		}, nil
	}

	attributeValues, err := attributevalue.MarshalMap(todo)
	if err != nil {
		return errorResponse, err
	}

	_, err = dynamodbClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      attributeValues,
	})
	if err != nil {
		return errorResponse, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 204,
	}, nil
}

func deleteTodo(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	id := request.PathParameters["id"]

	_, err := dynamodbClient.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		Key:       map[string]types.AttributeValue{"Id": &types.AttributeValueMemberS{Value: id}},
		TableName: &tableName,
	})
	if err != nil {
		return errorResponse, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 204,
	}, nil
}
