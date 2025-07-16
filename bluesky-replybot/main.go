package main

import (
	"context"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	genkitInstance, err := genkit.Init(context.Background(),
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
	)
	if err != nil {
		log.Fatalf("could not initialize Genkit: %v", err)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		authClient, auth, err := createAuthenticatedBlueskyClient()
		if err != nil {
			log.Fatal("Error creating authenticated Bluesky client:", err)
		}

		notifications, err := checkNotifications(authClient)
		if err != nil {
			log.Fatal("Error checking notifications:", err)
		}

		for _, notif := range notifications {
			processNotification(authClient, auth, genkitInstance, notif)
		}
	}
}

func processNotification(authClient *xrpc.Client, auth *atproto.ServerCreateSession_Output, genkitInstance *genkit.Genkit, notif *bsky.NotificationListNotifications_Notification) {
	var postText string
	feedPost, ok := notif.Record.Val.(*bsky.FeedPost)
	if !ok {
		log.Printf("Notification record is not a FeedPost: %v", notif.Record.Val)
		return
	}

	postText = feedPost.Text

	if postText == "" {
		return
	}

	cleanedText := strings.ReplaceAll(postText, "@llm.rasc.ch", "")
	cleanedText = strings.TrimSpace(cleanedText)

	if cleanedText == "" {
		return
	}

	aiResponse, err := generateGeminiResponse(genkitInstance, cleanedText)
	if err != nil {
		log.Printf("Failed to generate AI response: %v", err)
		return
	}

	err = sendReply(authClient, auth, notif, aiResponse)
	if err != nil {
		log.Printf("Failed to send reply: %v", err)
	}
}

func generateGeminiResponse(genkitInstance *genkit.Genkit, userMessage string) (string, error) {
	prompt := fmt.Sprintf(`You are a helpful AI assistant responding to a message on Bluesky (a social media platform similar to Twitter).
Please provide a thoughtful, engaging, and helpful response to the following message.
Keep your response concise and appropriate for social media (under 280 characters when possible).

User message: %s

Response:`, userMessage)

	resp, err := genkit.Generate(context.Background(), genkitInstance, ai.WithPrompt(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate response from Gemini: %w", err)
	}

	response := resp.Text()
	if response == "" {
		return "", fmt.Errorf("received empty response from Gemini")
	}

	return response, nil
}

func sendReply(authClient *xrpc.Client, auth *atproto.ServerCreateSession_Output, originalNotif *bsky.NotificationListNotifications_Notification, replyText string) error {

	replyRecord := bsky.FeedPost{
		Text:      replyText,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Reply: &bsky.FeedPost_ReplyRef{
			Root: &atproto.RepoStrongRef{
				Uri: originalNotif.Uri,
				Cid: originalNotif.Cid,
			},
			Parent: &atproto.RepoStrongRef{
				Uri: originalNotif.Uri,
				Cid: originalNotif.Cid,
			},
		},
	}

	encodedRecord := &util.LexiconTypeDecoder{Val: &replyRecord}

	_, err := atproto.RepoCreateRecord(
		context.Background(),
		authClient,
		&atproto.RepoCreateRecord_Input{
			Repo:       auth.Did,
			Collection: "app.bsky.feed.post",
			Record:     encodedRecord,
		},
	)

	return err
}

func checkNotifications(authClient *xrpc.Client) ([]*bsky.NotificationListNotifications_Notification, error) {
	limit := int64(5)
	reasons := []string{"mention"}
	cursor := ""
	var allUnreadNotifications []*bsky.NotificationListNotifications_Notification

	for {
		notificationsList, err := bsky.NotificationListNotifications(context.Background(), authClient, cursor, limit, false, reasons, "")
		if err != nil {
			return nil, fmt.Errorf("failed to list notifications: %w", err)
		}

		if len(notificationsList.Notifications) == 0 {
			break
		}

		hasReadMessages := false
		var unreadInBatch []*bsky.NotificationListNotifications_Notification

		for _, notif := range notificationsList.Notifications {
			if notif.IsRead {
				hasReadMessages = true
			} else {
				unreadInBatch = append(unreadInBatch, notif)
			}
		}

		allUnreadNotifications = append(allUnreadNotifications, unreadInBatch...)

		if hasReadMessages {
			break
		}

		if notificationsList.Cursor == nil {
			break
		}

		cursor = *notificationsList.Cursor
	}

	if len(allUnreadNotifications) > 0 {
		seenInput := &bsky.NotificationUpdateSeen_Input{
			SeenAt: time.Now().UTC().Format(time.RFC3339),
		}

		err := bsky.NotificationUpdateSeen(context.Background(), authClient, seenInput)
		if err != nil {
			return nil, fmt.Errorf("failed to mark notifications as seen: %w", err)
		}
	}

	return allUnreadNotifications, nil
}

func createAuthenticatedBlueskyClient() (*xrpc.Client, *atproto.ServerCreateSession_Output, error) {
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
		return nil, nil, fmt.Errorf("authentication failed: %w", err)
	}

	authClient := &xrpc.Client{
		Host: client.Host,
		Auth: &xrpc.AuthInfo{AccessJwt: auth.AccessJwt},
	}

	return authClient, auth, nil
}
