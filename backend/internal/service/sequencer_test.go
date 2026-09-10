package service

import (
	"testing"
	"time"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
)

func TestSequencerEmptyPlaylist(t *testing.T) {
	sequencer := NewSequencer()

	cycleStart := time.Unix(0, 0)
	now := cycleStart

	result := sequencer.GetCurrentPlayback(
		[]model.PlaylistItem{},
		cycleStart,
		now,
	)

	if result != nil {
		t.Fatal("expected nil for empty playlist")
	}
}

func TestSequencerFirstMedia(t *testing.T) {
	sequencer := NewSequencer()

	cycleStart := time.Unix(0, 0)
	now := cycleStart.Add(10 * time.Second)

	playlist := []model.PlaylistItem{
		{
			MediaID:  1,
			Position: 1,
			Duration: 30,
		},
		{
			MediaID:  2,
			Position: 2,
			Duration: 20,
		},
	}

	result := sequencer.GetCurrentPlayback(playlist, cycleStart, now)

	if result == nil {
		t.Fatal("expected playback state")
	}

	if result.MediaID != 1 {
		t.Fatalf("expected media 1, got %d", result.MediaID)
	}

	if result.Offset != 10 {
		t.Fatalf("expected offset 10, got %d", result.Offset)
	}
}

func TestSequencerMediaTransition(t *testing.T) {
	sequencer := NewSequencer()

	cycleStart := time.Unix(0, 0)
	now := cycleStart.Add(35 * time.Second)

	playlist := []model.PlaylistItem{
		{
			MediaID:  1,
			Position: 1,
			Duration: 30,
		},
		{
			MediaID:  2,
			Position: 2,
			Duration: 20,
		},
	}

	result := sequencer.GetCurrentPlayback(playlist, cycleStart, now)

	if result == nil {
		t.Fatal("expected playback state")
	}

	if result.MediaID != 2 {
		t.Fatalf("expected media 2, got %d", result.MediaID)
	}

	if result.Offset != 5 {
		t.Fatalf("expected offset 5, got %d", result.Offset)
	}
}

func TestSequencerLoopsPlaylist(t *testing.T) {
	sequencer := NewSequencer()

	cycleStart := time.Unix(0, 0)
	now := cycleStart.Add(55 * time.Second)

	playlist := []model.PlaylistItem{
		{
			MediaID:  1,
			Position: 1,
			Duration: 30,
		},
		{
			MediaID:  2,
			Position: 2,
			Duration: 20,
		},
	}

	result := sequencer.GetCurrentPlayback(playlist, cycleStart, now)

	if result == nil {
		t.Fatal("expected playback state")
	}

	if result.MediaID != 1 {
		t.Fatalf("expected media 1 after loop, got %d", result.MediaID)
	}

	if result.Offset != 5 {
		t.Fatalf("expected offset 5, got %d", result.Offset)
	}
}

func TestSequencerExactBoundary(t *testing.T) {
	sequencer := NewSequencer()

	cycleStart := time.Unix(0, 0)
	now := cycleStart.Add(30 * time.Second)

	playlist := []model.PlaylistItem{
		{
			MediaID:  1,
			Position: 1,
			Duration: 30,
		},
		{
			MediaID:  2,
			Position: 2,
			Duration: 20,
		},
	}

	result := sequencer.GetCurrentPlayback(playlist, cycleStart, now)

	if result == nil {
		t.Fatal("expected playback state")
	}

	if result.MediaID != 2 {
		t.Fatalf("expected media 2, got %d", result.MediaID)
	}

	if result.Offset != 0 {
		t.Fatalf("expected offset 0, got %d", result.Offset)
	}
}

func TestSequencerStartedAt(t *testing.T) {
	sequencer := NewSequencer()

	cycleStart := time.Unix(0, 0)
	now := cycleStart.Add(10 * time.Second)

	playlist := []model.PlaylistItem{
		{
			MediaID:  1,
			Position: 1,
			Duration: 30,
		},
	}

	result := sequencer.GetCurrentPlayback(playlist, cycleStart, now)

	if result == nil {
		t.Fatal("expected playback state")
	}

	expected := cycleStart

	if !result.StartedAt.Equal(expected) {
		t.Fatalf(
			"expected startedAt %v, got %v",
			expected,
			result.StartedAt,
		)
	}
}

func TestSequencerFiveHourBoundary(t *testing.T) {
	sequencer := NewSequencer()

	cycleStart := time.Unix(0, 0)

	playlist := []model.PlaylistItem{
		{
			MediaID:  1,
			Position: 1,
			Duration: 30,
		},
		{
			MediaID:  2,
			Position: 2,
			Duration: 20,
		},
	}

	// Exactly 5 hours later, the 5-hour cycle restarts.
	now := cycleStart.Add(5 * time.Hour)

	result := sequencer.GetCurrentPlayback(
		playlist,
		cycleStart,
		now,
	)

	if result == nil {
		t.Fatal("expected playback state")
	}

	if result.MediaID != 1 {
		t.Fatalf("expected media 1 after 5-hour cycle reset, got %d", result.MediaID)
	}

	if result.Offset != 0 {
		t.Fatalf("expected offset 0 after 5-hour cycle reset, got %d", result.Offset)
	}
}
