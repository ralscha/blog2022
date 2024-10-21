package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dcaraxes/gotenberg-go-client/v8"
)

func main() {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	client := &gotenberg.Client{Hostname: "http://localhost:3000", HTTPClient: httpClient}

	req := gotenberg.NewURLRequest("https://en.wikipedia.org/wiki/2024_World_Jigsaw_Puzzle_Championship")
	req.Format(gotenberg.PNG)
	req.OmitBackground()

	err := client.StoreScreenshot(req, "puzzle.png")
	if err != nil {
		log.Fatal(err)
	}
}
