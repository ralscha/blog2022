package internal

import (
	"fmt"
	"time"

	"google.golang.org/genai"
)

func GenRunWithRetries(userPrompt string, maxRetries int) (string, error) {
	if maxRetries <= 0 {
		return "", fmt.Errorf("maxRetries must be greater than 0, got %d", maxRetries)
	}

	var history []*genai.Content
	retryCount := 0

	for {
		fmt.Println("Generating code...")
		generatedCode := GenCode(userPrompt, history)

		retryCount++
		fmt.Println("Running code in Docker...")
		stdout, stderr, err := RunCodeInDocker(generatedCode, 5*time.Minute)

		if err != nil || stderr != "" {
			if retryCount >= maxRetries {
				if err != nil {
					return "", fmt.Errorf("error running code in Docker after %d retries: %v", maxRetries, err)
				}
				return "", fmt.Errorf("error running code in Docker after %d retries: %s", maxRetries, stderr)
			}

			errorMsg := ""
			if err != nil {
				errorMsg = fmt.Sprintf("Docker run failed with error: %v", err)
			} else {
				errorMsg = fmt.Sprintf("Docker run failed with stderr: %s", stderr)
			}

			fmt.Printf("Attempt %d failed: %s\n", retryCount, errorMsg)
			fmt.Println("Regenerating code with error context...")

			history = append(history,
				&genai.Content{
					Role:  "user",
					Parts: []*genai.Part{{Text: userPrompt}},
				},
				&genai.Content{
					Role:  "model",
					Parts: []*genai.Part{{Text: generatedCode}},
				})

			userPrompt = fmt.Sprintf(
				"The previous code execution failed with:\n\n%s\n\n"+
					"Please analyze the error and generate corrected Go code that:\n"+
					"1. Fixes the specific error\n"+
					"2. Maintains the original functionality", errorMsg)

			continue
		}

		return stdout, nil
	}
}
