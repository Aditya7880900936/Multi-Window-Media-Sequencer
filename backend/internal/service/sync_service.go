package service

import (
	"encoding/json"
	"time"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/repository"
	ws "github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/websocket"
	"gorm.io/gorm"
)

type SyncService struct {
	db  *gorm.DB
	hub *ws.Hub
}

func NewSyncService(db *gorm.DB, hub *ws.Hub) *SyncService {
	return &SyncService{
		db:  db,
		hub: hub,
	}
}

type SyncMessage struct {
	Type      string    `json:"type"`
	MediaID   uint      `json:"mediaId"`
	StartedAt time.Time `json:"startedAt"`
	Duration  int       `json:"duration"`
}

func (s *SyncService) Sync(mediaID uint, duration int) error {
	// Verify media exists.
	mediaRepo := repository.NewMediaRepository(s.db)

	_, err := mediaRepo.GetByID(mediaID)
	if err != nil {
		return err
	}

	startedAt := time.Now().UTC()

	event := &model.SyncEvent{
		MediaID:   mediaID,
		StartedAt: startedAt,
		Duration:  duration,
	}

	if err := s.db.Create(event).Error; err != nil {
		return err
	}

	message := SyncMessage{
		Type:      "SYNC_PLAY",
		MediaID:   mediaID,
		StartedAt: startedAt,
		Duration:  duration,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	s.hub.Broadcast <- data

	return nil
}