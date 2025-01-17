package main

import (
	"context"
	"fmt"
	"godockerllm/internal"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"detailedDescription": {
				Type:        genai.TypeString,
				Description: "The detailed description about a program the gen_code_tool should generate",
			},
		},
		Required: []string{"detailedDescription"},
	}

	generalPurposeTool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name: "gen_code_tool",
			Description: `This is a general-purpose tool that can be used to generate code for various tasks,
including searching for information about competitions and events and accessing public APIs and accessing real-time information
and has access to the Internet. The tool generates Go code for the given detailed description and runs it in a Docker container.
It returns the output of the code execution. `,
			Parameters: schema,
		}},
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	model.Tools = []*genai.Tool{generalPurposeTool}
	session := model.StartChat()
	// userPrompt := "What is the current date and time"
	userPrompt := "Is the number 1201281 prime? List all divisors of this number?"
	res, err := session.SendMessage(ctx, genai.Text(userPrompt))
	if err != nil {
		log.Fatalf("session.SendMessage: %v", err)
	}

	part := res.Candidates[0].Content.Parts[0]
	funcall, ok := part.(genai.FunctionCall)
	if !ok {
		fmt.Println("No Function Calling")
		fmt.Println(part)
		return
	}

	if funcall.Name != generalPurposeTool.FunctionDeclarations[0].Name {
		log.Fatalf("expected %q, got %q", generalPurposeTool.FunctionDeclarations[0].Name, funcall.Name)
	}

	description, ok := funcall.Args["detailedDescription"].(string)
	if !ok {
		log.Fatalf("expected string: %v", funcall.Args["detailedDescription"])
	}

	result, err := internal.GenRunWithRetries(description, 3)
	if err != nil {
		log.Fatal(err)
	}

	res, err = session.SendMessage(ctx, genai.FunctionResponse{
		Name: generalPurposeTool.FunctionDeclarations[0].Name,
		Response: map[string]any{
			"result": result,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, cand := range res.Candidates {
		if cand.Content != nil {
			for _, part2 := range cand.Content.Parts {
				fmt.Println(part2)
			}
		}
	}
}
