package service

import (
	"time"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/repository"
)

type WindowService struct {
	repo        *repository.WindowRepository
	sequencer   *Sequencer
	cycleRepo   *repository.PlaybackCycleRepository
	syncService *SyncService
}

func NewWindowService(
	repo *repository.WindowRepository,
	sequencer *Sequencer,
	cycleRepo *repository.PlaybackCycleRepository,
	syncService *SyncService,
) *WindowService {
	return &WindowService{
		repo:        repo,
		sequencer:   sequencer,
		cycleRepo:   cycleRepo,
		syncService: syncService,
	}
}

func (s *WindowService) GetAll() ([]model.Window, error) {
	return s.repo.GetAll()
}

func (s *WindowService) GetByID(id uint) (*model.Window, error) {
	return s.repo.GetByID(id)
}

func (s *WindowService) GetCurrentPlayback(
	id uint,
	now time.Time,
) (*PlaybackState, error) {

	window, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Sync takes priority over the normal playlist.
	activeSync, err := s.syncService.GetActiveSync(now)
	if err != nil {
		return nil, err
	}

	if activeSync != nil {
		return &PlaybackState{
			MediaID:   activeSync.MediaID,
			Position:  0,
			Offset:    activeSync.Offset,
			StartedAt: activeSync.StartedAt,
		}, nil
	}

	cycle, err := s.cycleRepo.GetOrCreate()
	if err != nil {
		return nil, err
	}

	return s.sequencer.GetCurrentPlayback(
		window.Playlist,
		cycle.StartedAt,
		now,
	), nil
}
