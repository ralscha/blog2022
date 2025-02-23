package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/joho/godotenv"
)

func main() {
	client := &xrpc.Client{
		Host: "https://bsky.social",
	}

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	did := os.Getenv("BLUESKY_IDENTIFIER")
	appPassword := os.Getenv("BLUESKY_PASSWORD")

	if did == "" || appPassword == "" {
		log.Fatal("Missing required environment variables: BLUESKY_IDENTIFIER and BLUESKY_PASSWORD")
	}

	auth, err := atproto.ServerCreateSession(
		context.Background(),
		client,
		&atproto.ServerCreateSession_Input{
			Identifier: did,
			Password:   appPassword,
		},
	)
	if err != nil {
		log.Fatal("Failed to authenticate:", err)
	}

	authClient := xrpc.Client{
		Host: client.Host,
		Auth: &xrpc.AuthInfo{AccessJwt: auth.AccessJwt},
	}

	currentTime := time.Now()
	post := &bsky.FeedPost{
		Text:      fmt.Sprintf("👋 Hello World! Posted from my Go program at %s", currentTime.Format("15:04:05")),
		CreatedAt: currentTime.Format(time.RFC3339),
	}

	_, err = atproto.RepoCreateRecord(
		context.Background(),
		&authClient,
		&atproto.RepoCreateRecord_Input{
			Repo:       auth.Did,
			Collection: "app.bsky.feed.post",
			Record:     &util.LexiconTypeDecoder{Val: post},
		},
	)
	if err != nil {
		log.Fatal("Failed to create post:", err)
	}

	fmt.Println("Successfully posted 'Hello World!' to Blue Sky!")
}
