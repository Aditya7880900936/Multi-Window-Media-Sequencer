package service

import (
	"time"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
)

type PlaybackState struct {
	MediaID   uint      `json:"mediaId"`
	Position  int       `json:"position"`
	Offset    int       `json:"offset"`
	StartedAt time.Time `json:"startedAt"`
}
type Sequencer struct {
	CycleDuration time.Duration
}

func NewSequencer() *Sequencer {
	return &Sequencer{
		CycleDuration: 5 * time.Hour,
	}
}

func (s *Sequencer) GetCurrentPlayback(
	playlist []model.PlaylistItem,
	cycleStart time.Time,
	now time.Time,
) *PlaybackState {
	if len(playlist) == 0 {
		return nil
	}

	var playlistDuration time.Duration

	for _, item := range playlist {
		if item.Duration <= 0 {
			continue
		}

		playlistDuration += time.Duration(item.Duration) * time.Second
	}

	if playlistDuration <= 0 {
		return nil
	}

	elapsed := now.Sub(cycleStart)

	if elapsed < 0 {
		elapsed = 0
	}

	// Continuously loop through the configured playlist.
	positionInPlaylist := elapsed % playlistDuration

	current := positionInPlaylist

	for _, item := range playlist {
		if item.Duration <= 0 {
			continue
		}

		duration := time.Duration(item.Duration) * time.Second

		if current < duration {
			offset := int(current.Seconds())

			return &PlaybackState{
				MediaID:   item.MediaID,
				Position:  item.Position,
				Offset:    offset,
				StartedAt: now.Add(-current),
			}
		}

		current -= duration
	}

	return nil
}
