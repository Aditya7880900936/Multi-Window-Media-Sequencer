package service

import (
	"encoding/json"
	"fmt"
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

type SyncPlaybackState struct {
	MediaID   uint
	Offset    int
	StartedAt time.Time
	Duration  int
}

func (s *SyncService) Sync(mediaID uint, duration int) error {
	mediaRepo := repository.NewMediaRepository(s.db)

	_, err := mediaRepo.GetByID(mediaID)
	if err != nil {
		return err
	}

	if duration <= 0 {
		return fmt.Errorf("duration must be greater than zero")
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

	// Replace any currently active sync.
	activeSync := model.ActiveSync{
		ID:        1,
		MediaID:   mediaID,
		StartedAt: startedAt,
		Duration:  duration,
	}

	if err := s.db.Save(&activeSync).Error; err != nil {
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

func (s *SyncService) GetActiveSync(
	now time.Time,
) (*SyncPlaybackState, error) {
	repo := repository.NewActiveSyncRepository(s.db)

	activeSync, err := repo.Get()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	elapsed := now.Sub(activeSync.StartedAt)

	if elapsed < 0 {
		elapsed = 0
	}

	if elapsed >= time.Duration(activeSync.Duration)*time.Second {
		if err := repo.Clear(); err != nil {
			return nil, err
		}

		return nil, nil
	}

	return &SyncPlaybackState{
		MediaID:   activeSync.MediaID,
		Offset:    int(elapsed.Seconds()),
		StartedAt: activeSync.StartedAt,
		Duration:  activeSync.Duration,
	}, nil
}
