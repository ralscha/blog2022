package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

func TestStateSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	state, err := loadState(path)
	if err != nil {
		t.Fatal(err)
	}
	state.Pending["at://pending"] = &bsky.NotificationListNotifications_Notification{
		Uri: "at://pending", Record: &util.LexiconTypeDecoder{Val: &bsky.FeedPost{Text: "hello"}},
	}
	state.Completed["at://done"] = true
	if err := saveState(path, state); err != nil {
		t.Fatal(err)
	}
	state, err = loadState(path)
	if err != nil {
		t.Fatal(err)
	}
	if !state.Completed["at://done"] || state.Pending["at://pending"].Record.Val.(*bsky.FeedPost).Text != "hello" {
		t.Fatal("pending or completed work was lost")
	}
	delete(state.Pending, "at://pending")
	if err := saveState(path, state); err != nil {
		t.Fatal(err)
	}
	state, err = loadState(path)
	if err != nil || len(state.Pending) != 0 {
		t.Fatalf("state replacement failed: %v", err)
	}
}

func TestTokenRefreshWindow(t *testing.T) {
	now := time.Unix(1_000, 0)
	for _, expiry := range []int64{999, 1060, 1061, 7200} {
		payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"exp":%d}`, expiry)))
		if got := tokenNeedsRefresh("header."+payload+".signature", now); got != (expiry <= 1060) {
			t.Fatalf("unexpected refresh decision for expiry %d", expiry)
		}
	}
	if !tokenNeedsRefresh("malformed", now) {
		t.Fatal("malformed token was accepted")
	}
}

func TestReplyPreservesThreadRootAndStableKey(t *testing.T) {
	root := &atproto.RepoStrongRef{Uri: "at://root", Cid: "root-cid"}
	notif := &bsky.NotificationListNotifications_Notification{
		Uri: "at://parent", Cid: "parent-cid",
		Record: &util.LexiconTypeDecoder{Val: &bsky.FeedPost{Reply: &bsky.FeedPost_ReplyRef{Root: root}}},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Rkey   string        `json:"rkey"`
			Record bsky.FeedPost `json:"record"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body.Rkey != replyKey(notif.Uri) || body.Record.Reply.Root.Uri != root.Uri || body.Record.Reply.Parent.Uri != notif.Uri {
			t.Errorf("incorrect reply references or key: %+v", body)
		}
		fmt.Fprint(w, `{"uri":"at://reply","cid":"reply-cid"}`)
	}))
	defer server.Close()
	if err := sendReply(&xrpc.Client{Host: server.URL}, &atproto.ServerCreateSession_Output{Did: "did:example:bot"}, notif, "reply"); err != nil {
		t.Fatal(err)
	}
}
