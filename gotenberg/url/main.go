package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dcaraxes/gotenberg-go-client/v8"
)

func main() {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	client := &gotenberg.Client{Hostname: "http://localhost:3000", HTTPClient: httpClient}

	req := gotenberg.NewURLRequest("https://xkcd.com/")
	req.Margins(gotenberg.NoMargins)
	req.Scale(0.9)
	req.SinglePage()

	response, err := client.Post(req)
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
