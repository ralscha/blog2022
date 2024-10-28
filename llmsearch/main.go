package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/ollama/ollama/api"
)

type Result struct {
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	Content       string   `json:"content"`
	Engine        string   `json:"engine"`
	ParsedURL     []string `json:"parsed_url"`
	Template      string   `json:"template"`
	Engines       []string `json:"engines"`
	Positions     []int    `json:"positions"`
	Thumbnail     string   `json:"thumbnail"`
	PublishedDate *string  `json:"publishedDate"`
	Score         float64  `json:"score"`
	Category      string   `json:"category"`
}

type SearchResponse struct {
	Query           string   `json:"query"`
	NumberOfResults int      `json:"number_of_results"`
	Results         []Result `json:"results"`
}

const (
	maxResults    = 3
	searchBaseURL = "http://localhost:8080/search"
	ollamaBaseURL = "http://localhost:11434"
	modelName     = "qwen2.5:7b"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Minute,
}

func main() {
	ollamaURL, err := url.Parse(ollamaBaseURL)
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(ollamaURL, httpClient)

	query := "Who won the Puzzle World Championship 2024?"
	searchQuery := getSearchQuery(client, query)
	fmt.Println("Query:", searchQuery)

	searchResponse, err := webSearch(searchQuery)
	if err != nil {
		log.Fatal(err)
	}

	searchContext := buildSearchContext(searchResponse.Results)
	answer := getAnswer(client, query, searchContext)
	fmt.Println(answer)
}

func getSearchQuery(client *api.Client, query string) string {
	messages := []api.Message{
		{
			Role:    "system",
			Content: "You are a professional web searcher.",
		},
		{
			Role:    "user",
			Content: "Reformulate the following user prompt into a search query and return it.Nothing else.\n\n" + query,
		},
	}

	response := executeChat(client, messages)
	return strings.Trim(response, "\"")
}

func getAnswer(client *api.Client, query, context string) string {
	messages := []api.Message{
		{
			Role: "user",
			Content: fmt.Sprintf("%s\n\nOnly return answer based on the context. "+
				"If you don't know return I don't know:\n###%s\n###", query, context),
		},
	}

	return executeChat(client, messages)
}

func executeChat(client *api.Client, messages []api.Message) string {
	req := &api.ChatRequest{
		Model:    modelName,
		Messages: messages,
	}

	var wg sync.WaitGroup
	var sb strings.Builder
	var response string

	respFunc := func(resp api.ChatResponse) error {
		sb.WriteString(resp.Message.Content)
		if resp.Done {
			response = sb.String()
			wg.Done()
		}
		return nil
	}

	wg.Add(1)
	ctx := context.Background()
	if err := client.Chat(ctx, req, respFunc); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	return response
}

func buildSearchContext(results []Result) string {
	var wg sync.WaitGroup
	resultChan := make(chan string, maxResults)

	for i, result := range results {
		if i >= maxResults {
			break
		}

		wg.Add(1)
		go func(result Result) {
			defer wg.Done()
			content, err := fetchTextContent(result.URL)
			if err != nil {
				log.Printf("Error fetching content for URL %s: %v", result.URL, err)
				return
			}
			resultChan <- fmt.Sprintf("%s\n%s\n\n", result.URL, content)
		}(result)
	}

	wg.Wait()
	close(resultChan)

	var contextSB strings.Builder
	for res := range resultChan {
		contextSB.WriteString(res)
	}
	return contextSB.String()
}

func fetchTextContent(url string) (string, error) {
	fmt.Println("Fetching content for URL:", url)
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return "", err
	}
	return article.TextContent, nil
}

func webSearch(query string) (*SearchResponse, error) {
	encodedQuery := url.QueryEscape(query)
	requestURL := fmt.Sprintf("%s?q=%s&format=json", searchBaseURL, encodedQuery)

	response, err := httpClient.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var searchResponse SearchResponse
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	return &searchResponse, nil
}
