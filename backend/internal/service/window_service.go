package service

import (
	"time"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/repository"
)

type WindowService struct {
	repo      *repository.WindowRepository
	sequencer *Sequencer
	cycleRepo *repository.PlaybackCycleRepository
}

func NewWindowService(
	repo *repository.WindowRepository,
	sequencer *Sequencer,
	cycleRepo *repository.PlaybackCycleRepository,
) *WindowService {
	return &WindowService{
		repo:      repo,
		sequencer: sequencer,
		cycleRepo: cycleRepo,
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
