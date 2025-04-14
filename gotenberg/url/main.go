package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
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

	req := gotenberg.NewURLRequest("https://xkcd.com/")
	req.Margins(gotenberg.NoMargins)
	req.Scale(0.9)
	req.SinglePage()

	ctx := context.Background()
	response, err := client.Send(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	file, err := os.Create("xkcd.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

}
