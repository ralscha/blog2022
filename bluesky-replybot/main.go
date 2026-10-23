package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	genkitInstance := genkit.Init(context.Background(),
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
	)

	authClient, auth, err := createAuthenticatedBlueskyClient()
	if err != nil {
		log.Fatal(err)
	}
	statePath := os.Getenv("BLUESKY_STATE_FILE")
	if statePath == "" {
		statePath = "replybot-state.json"
	}
	state, err := loadState(statePath)
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		if err := refreshSession(authClient, auth); err != nil {
			log.Printf("Session refresh failed: %v", err)
		} else if err := poll(authClient, auth, genkitInstance, state, statePath); err != nil {
			log.Printf("Polling failed; pending work will be retried: %v", err)
		}
		<-ticker.C
	}
}

func poll(authClient *xrpc.Client, auth *atproto.ServerCreateSession_Output, genkitInstance *genkit.Genkit, state *botState, statePath string) error {
	notifications, err := checkNotifications(authClient)
	if err != nil {
		return err
	}
	var seenAt string
	for _, notif := range notifications {
		if notif.IndexedAt > seenAt {
			seenAt = notif.IndexedAt
		}
		if !state.Completed[notif.Uri] {
			state.Pending[notif.Uri] = notif
		}
	}
	if err := saveState(statePath, state); err != nil {
		return err
	}
	for uri, notif := range state.Pending {
		if err := processNotification(authClient, auth, genkitInstance, notif); err != nil {
			return err
		}
		state.Completed[uri] = true
		delete(state.Pending, uri)
		if err := saveState(statePath, state); err != nil {
			return err
		}
	}
	if seenAt != "" {
		return bsky.NotificationUpdateSeen(context.Background(), authClient, &bsky.NotificationUpdateSeen_Input{SeenAt: seenAt})
	}
	return nil
}

func processNotification(authClient *xrpc.Client, auth *atproto.ServerCreateSession_Output, genkitInstance *genkit.Genkit, notif *bsky.NotificationListNotifications_Notification) error {
	feedPost, ok := notif.Record.Val.(*bsky.FeedPost)
	if !ok {
		return nil
	}
	cleanedText := strings.TrimSpace(strings.ReplaceAll(feedPost.Text, "@llm.rasc.ch", ""))
	if cleanedText == "" {
		return nil
	}
	// A stable record key also covers a crash after posting but before saving local state.
	_, err := atproto.RepoGetRecord(context.Background(), authClient, "", "app.bsky.feed.post", auth.Did, replyKey(notif.Uri))
	if err == nil {
		return nil
	}
	var rpcError *xrpc.XRPCError
	if !errors.As(err, &rpcError) || rpcError.ErrStr != "RecordNotFound" {
		return err
	}
	aiResponse, err := generateGeminiResponse(genkitInstance, cleanedText)
	if err != nil {
		return err
	}
	return sendReply(authClient, auth, notif, aiResponse)
}

func replyKey(uri string) string {
	return fmt.Sprintf("reply-%x", sha256.Sum256([]byte(uri)))
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

	if original, ok := originalNotif.Record.Val.(*bsky.FeedPost); ok && original.Reply != nil {
		replyRecord.Reply.Root = original.Reply.Root
	}
	rkey := replyKey(originalNotif.Uri)
	encodedRecord := &util.LexiconTypeDecoder{Val: &replyRecord}

	_, err := atproto.RepoCreateRecord(
		context.Background(),
		authClient,
		&atproto.RepoCreateRecord_Input{
			Repo:       auth.Did,
			Collection: "app.bsky.feed.post",
			Record:     encodedRecord,
			Rkey:       &rkey,
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
		Auth: &xrpc.AuthInfo{AccessJwt: auth.AccessJwt, RefreshJwt: auth.RefreshJwt, Did: auth.Did, Handle: auth.Handle},
	}

	return authClient, auth, nil
}
