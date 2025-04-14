package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/starwalkn/gotenberg-go-client/v8"
)

func main() {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	client, err := gotenberg.NewClient("http://localhost:3000", httpClient)
	if err != nil {
		log.Fatal(err)
	}

	req := gotenberg.NewURLRequest("https://en.wikipedia.org/wiki/2024_World_Jigsaw_Puzzle_Championship")
	req.Format(gotenberg.PNG)
	req.OmitBackground()

	ctx := context.Background()
	err = client.StoreScreenshot(ctx, req, "puzzle.png")
	if err != nil {
		log.Fatal(err)
	}
}
