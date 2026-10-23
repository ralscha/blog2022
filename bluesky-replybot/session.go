package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

func refreshSession(client *xrpc.Client, auth *atproto.ServerCreateSession_Output) error {
	if !tokenNeedsRefresh(auth.AccessJwt, time.Now()) {
		return nil
	}
	refreshClient := *client
	refreshClient.Auth = &xrpc.AuthInfo{AccessJwt: auth.RefreshJwt}
	refreshed, err := atproto.ServerRefreshSession(context.Background(), &refreshClient)
	if err != nil {
		var rpcError *xrpc.XRPCError
		if !errors.As(err, &rpcError) || (rpcError.ErrStr != "ExpiredToken" && rpcError.ErrStr != "InvalidToken") {
			return err
		}
		newClient, newAuth, loginErr := createAuthenticatedBlueskyClient()
		if loginErr != nil {
			return loginErr
		}
		*client, *auth = *newClient, *newAuth
		return nil
	}
	auth.AccessJwt, auth.RefreshJwt = refreshed.AccessJwt, refreshed.RefreshJwt
	client.Auth.AccessJwt, client.Auth.RefreshJwt = refreshed.AccessJwt, refreshed.RefreshJwt
	return nil
}

func tokenNeedsRefresh(token string, now time.Time) bool {
	// Decode only to schedule refresh; the PDS authenticates every API request.
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return true
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return true
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if json.Unmarshal(payload, &claims) != nil {
		return true
	}
	return !time.Unix(claims.Exp, 0).After(now.Add(time.Minute))
}
