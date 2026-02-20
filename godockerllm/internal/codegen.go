package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/genai"
)

func GenCode(userPrompt string, history []*genai.Content) string {
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
		log.Fatalf("Error creating client: %v", err)
	}

	systemPrompt := `
Write a Go program that meets the following requirements:
- The code should follow best practices.
- Ensure the code is efficient and optimized.
- Handle errors gracefully and provide meaningful error messages.
- Use idiomatic Go constructs and conventions.
- Do not comment the code
- Only return the Go code. Nothing else should be returned.
`
	var parts []*genai.Part
	for _, h := range history {
		parts = append(parts, h.Parts...)
	}
	parts = append(parts, &genai.Part{Text: userPrompt})

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", []*genai.Content{{Parts: parts}}, &genai.GenerateContentConfig{
		Temperature:       new(float32(0.7)),
		MaxOutputTokens:   8192,
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemPrompt}}},
	})
	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		var llmOutput strings.Builder
		for _, part := range result.Candidates[0].Content.Parts {
			llmOutput.WriteString(fmt.Sprintf("%v", part.Text))
		}

		return cleanup(llmOutput.String())
	}

	log.Fatalln("No response received from LLM")
	return ""
}

func cleanup(code string) string {
	updatedCode := strings.TrimPrefix(code, "```go")
	updatedCode = strings.TrimSuffix(updatedCode, "\n")
	updatedCode = strings.TrimSuffix(updatedCode, "```")
	return updatedCode
}
