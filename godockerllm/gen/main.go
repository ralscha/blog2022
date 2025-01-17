package main

import (
	"fmt"
	"godockerllm/internal"
)

func main() {
	userPrompt := "Write a Go program that prints the first 10 Fibonacci numbers."
	generatedCode := internal.GenCode(userPrompt, nil)
	fmt.Println(generatedCode)
}
