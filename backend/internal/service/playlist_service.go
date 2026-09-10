package service

import (
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/repository"
)

type PlaylistService struct {
	repo *repository.PlaylistRepository
}

func NewPlaylistService(repo *repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{repo: repo}
}

func (s *PlaylistService) GetByWindowID(windowID uint) ([]model.PlaylistItem, error) {
	return s.repo.GetByWindowID(windowID)
}

func (s *PlaylistService) Create(item *model.PlaylistItem) error {
	return s.repo.Create(item)
}

func (s *PlaylistService) Delete(windowID, itemID uint) error {
	return s.repo.Delete(windowID, itemID)
}