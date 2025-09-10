package main

import (
	"context"
	"fmt"
	"godockerllm/internal"
	"log"
	"os"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

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

	userPrompt := "Is the number 1201281 prime? List all divisors of this number?"

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: userPrompt}}, Role: "user"},
	}
	config := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(float32(0.7)),
		Tools:       []*genai.Tool{generalPurposeTool},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}

	var funcCall *genai.FunctionCall
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if part.FunctionCall != nil {
					funcCall = part.FunctionCall
					break
				}
			}
		}
	}
	if funcCall == nil {
		fmt.Println("No Function Calling")
		fmt.Println(resp.Text())
		return
	}

	description, ok := funcCall.Args["detailedDescription"].(string)
	if !ok {
		log.Fatalf("expected string: %v", funcCall.Args["detailedDescription"])
	}
	fmt.Println("Generated code description:", description)

	genresult, err := internal.GenRunWithRetries(description, 3)
	if err != nil {
		log.Fatal(err)
	}

	contents = append(contents, &genai.Content{
		Parts: []*genai.Part{{FunctionCall: funcCall}},
		Role:  "model",
	})
	contents = append(contents, &genai.Content{
		Parts: []*genai.Part{{FunctionResponse: &genai.FunctionResponse{
			Name: funcCall.Name,
			Response: map[string]any{
				"result": genresult,
			},
		}}},
		Role: "function",
	})

	resp, err = client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Text())
}
