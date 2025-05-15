package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/cockroachdb/pebble"
	"github.com/joho/godotenv"
)

type Earthquake struct {
	Time   string
	Mag    float64
	Status string
	ID     string
	Place  string
	Type   string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	resp, err := http.Get("https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/4.5_hour.csv")
	if err != nil {
		log.Fatal("Failed to download CSV:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	headers, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read CSV header:", err)
	}

	var earthquakes []Earthquake
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Failed to read CSV record:", err)
		}

		quakeMap := make(map[string]string)
		for i, h := range headers {
			quakeMap[h] = record[i]
		}

		mag, err := strconv.ParseFloat(quakeMap["mag"], 64)
		if err != nil {
			log.Printf("Skipping invalid magnitude '%s' for ID %s", quakeMap["mag"], quakeMap["id"])
			continue
		}

		earthquakes = append(earthquakes, Earthquake{
			Time:   quakeMap["time"],
			Mag:    mag,
			Status: quakeMap["status"],
			ID:     quakeMap["id"],
			Place:  quakeMap["place"],
			Type:   quakeMap["type"],
		})
	}

	var filtered []Earthquake
	for _, q := range earthquakes {
		if q.Mag >= 5 && q.Status == "reviewed" {
			filtered = append(filtered, q)
		}
	}

	db, err := pebble.Open("quake-db", &pebble.Options{})
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	for _, q := range filtered {
		key := []byte(q.ID)
		_, closer, err := db.Get(key)
		if err == nil {
			closer.Close()
			continue
		}
		if !errors.Is(err, pebble.ErrNotFound) {
			log.Printf("Database error for ID %s: %v", q.ID, err)
			continue
		}

		t, err := time.Parse("2006-01-02T15:04:05Z", q.Time)
		if err != nil {
			log.Printf("Failed to parse time for ID %s: %v", q.ID, err)
			continue
		}

		isoTimestamp := t.Format("2006-01-02 15:04:05 UTC")

		fullURL := fmt.Sprintf("https://earthquake.usgs.gov/earthquakes/eventpage/%s/executive", q.ID)
		shortURL := fmt.Sprintf("earthquake.usgs.gov/%s", q.ID)

		msg := fmt.Sprintf("%.1f magnitude %s #%s\n%s\n%s\n\n%s",
			q.Mag, earthquakeTypeByMagnitude(q.Mag), q.Type, isoTimestamp, q.Place, shortURL)

		if err := postToBluesky(msg, q.Type, fullURL, shortURL); err != nil {
			log.Printf("Failed to post for ID %s: %v", q.ID, err)
			continue
		}

		if err := db.Set(key, []byte("posted"), &pebble.WriteOptions{}); err != nil {
			log.Printf("Failed to store ID %s: %v", q.ID, err)
		}
	}
	_ = db.Flush()
}

func postToBluesky(text string, earthquakeType string, fullURL string, shortURL string) error {
	client := &xrpc.Client{Host: "https://me.rasc.ch"}

	auth, err := atproto.ServerCreateSession(
		context.Background(),
		client,
		&atproto.ServerCreateSession_Input{
			Identifier: os.Getenv("BLUESKY_IDENTIFIER"),
			Password:   os.Getenv("BLUESKY_PASSWORD"),
		},
	)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	authClient := xrpc.Client{
		Host: client.Host,
		Auth: &xrpc.AuthInfo{AccessJwt: auth.AccessJwt},
	}

	linkStartPos := strings.Index(text, shortURL)
	linkEndPos := linkStartPos + len(shortURL)

	linkFacet := &bsky.RichtextFacet{
		Features: []*bsky.RichtextFacet_Features_Elem{
			{
				RichtextFacet_Link: &bsky.RichtextFacet_Link{
					Uri: fullURL,
				},
			},
		},
		Index: &bsky.RichtextFacet_ByteSlice{
			ByteEnd:   int64(linkEndPos),
			ByteStart: int64(linkStartPos),
		},
	}

	facets := []*bsky.RichtextFacet{linkFacet}

	tagStartPos := strings.Index(text, "#"+earthquakeType)
	if tagStartPos != -1 {
		tagEndPos := tagStartPos + len(earthquakeType) + 1

		tagFacet := &bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Tag: &bsky.RichtextFacet_Tag{
						Tag: earthquakeType,
					},
				},
			},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteEnd:   int64(tagEndPos),
				ByteStart: int64(tagStartPos),
			},
		}
		facets = append(facets, tagFacet)
	}

	post := &bsky.FeedPost{
		Text:      text,
		Langs:     []string{"en"},
		CreatedAt: time.Now().Format(time.RFC3339),
		Facets:    facets,
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
	return err
}

func earthquakeTypeByMagnitude(mag float64) string {
	switch {
	case mag >= 8.0:
		return "great"
	case mag >= 7.0:
		return "major"
	case mag >= 6.0:
		return "strong"
	case mag >= 5.0:
		return "moderate"
	default:
		return ""
	}
}
