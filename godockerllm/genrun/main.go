package main

import (
	"fmt"
	"godockerllm/internal"
	"log"
)

func main() {
	userPrompt := "Write a Go program that prints the first 10 Fibonacci numbers as a comma-separated string."
	result, err := internal.GenRunWithRetries(userPrompt, 3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Output:")
	fmt.Println(result)

	// Wikipedia Search example
	userPrompt = "Search Wikipedia for the definition of the word 'gopher'."
	result, err = internal.GenRunWithRetries(userPrompt, 3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Output:")
	fmt.Println(result)

}
