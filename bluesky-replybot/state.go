package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/bluesky-social/indigo/api/bsky"
)

type botState struct {
	Pending   map[string]*bsky.NotificationListNotifications_Notification `json:"pending"`
	Completed map[string]bool                                             `json:"completed"`
}

func loadState(path string) (*botState, error) {
	state := &botState{}
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err == nil {
		if err := json.Unmarshal(data, state); err != nil {
			return nil, err
		}
	}
	if state.Pending == nil {
		state.Pending = make(map[string]*bsky.NotificationListNotifications_Notification)
	}
	if state.Completed == nil {
		state.Completed = make(map[string]bool)
	}
	return state, nil
}

func saveState(path string, state *botState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".replybot-state-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}
