package internal

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

func GenCode(userPrompt string, history []*genai.Content) string {
	ctx := context.Background()
	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash-exp")
	model.SetTemperature(0.7)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	systemPrompt := `
Write a Go program that meets the following requirements:
- The code should follow best practices.
- Ensure the code is efficient and optimized.
- Handle errors gracefully and provide meaningful error messages.
- Use idiomatic Go constructs and conventions.
- Do not comment the code
- Only return the Go code. Nothing else should be returned.
`
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	session := model.StartChat()
	session.History = history
	resp, err := session.SendMessage(ctx, genai.Text(userPrompt))
	if err != nil {
		log.Fatalf("Error sending message to LLM: %v", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		llmOutput := ""
		for _, part := range resp.Candidates[0].Content.Parts {
			llmOutput += fmt.Sprintf("%v", part)
		}

		return cleanup(llmOutput)
	}

	log.Fatalln("No response received from LLM")
	return ""
}

func cleanup(code string) string {
	updatedCode := strings.TrimPrefix(code, "```go")
	return strings.TrimSuffix(updatedCode, "```\n")
}
