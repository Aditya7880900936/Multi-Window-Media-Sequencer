package service

import (
	"encoding/json"

	ws "github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/websocket"
)

type PlaylistUpdateMessage struct {
	Type     string `json:"type"`
	WindowID uint   `json:"windowId"`
}

func BroadcastPlaylistUpdate(
	hub *ws.Hub,
	windowID uint,
) error {
	message := PlaylistUpdateMessage{
		Type:     "PLAYLIST_UPDATED",
		WindowID: windowID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	hub.Broadcast <- data

	return nil
}
