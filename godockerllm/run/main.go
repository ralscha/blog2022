package main

import (
	"fmt"
	"godockerllm/internal"
	"log"
	"time"
)

func main() {
	goProgram := `
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
`
	stdout, stderr, err := internal.RunCodeInDocker(goProgram, 5*time.Minute)
	if err != nil {
		log.Fatalf("Error running code in Docker: %v", err)
	}

	if stderr != "" {
		fmt.Println("Error:")
		fmt.Println(stderr)
	} else {
		fmt.Println("Output:")
		fmt.Println(stdout)
	}
}
