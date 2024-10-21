package main

import (
	"github.com/dcaraxes/gotenberg-go-client/v8"
	"io"
	"log"
	"net/http"
	"time"
)

const html = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Gopher</title>
    <link rel="stylesheet" href="style.css">
  </head>
  <body>
    <h1>Gopher</h1>
    <img src="gopher.png" width="100" />
  </body>
</html>
`

const css = `
body {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
  background-color: green;
}
h1 {
  color: black;
  font-size: 6em;
}
`

func main() {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	gopherURL := "https://raw.githubusercontent.com/golang-samples/gopher-vector/refs/heads/master/gopher.png"

	gopherResp, err := httpClient.Get(gopherURL)
	if err != nil {
		log.Fatal(err)
	}
	defer gopherResp.Body.Close()

	gopherBytes, err := io.ReadAll(gopherResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	client := &gotenberg.Client{Hostname: "http://localhost:3000", HTTPClient: httpClient}

	index, err := gotenberg.NewDocumentFromString("index.html", html)
	if err != nil {
		log.Fatal(err)
	}

	style, err := gotenberg.NewDocumentFromString("style.css", css)
	if err != nil {
		log.Fatal(err)
	}

	gopher, err := gotenberg.NewDocumentFromBytes("gopher.png", gopherBytes)
	if err != nil {
		log.Fatal(err)
	}

	req := gotenberg.NewHTMLRequest(index)
	req.Assets(style, gopher)
	req.PaperSize(gotenberg.A4)
	req.Margins(gotenberg.NoMargins)

	err = client.Store(req, "my.pdf")
	if err != nil {
		log.Fatal(err)
	}
}
