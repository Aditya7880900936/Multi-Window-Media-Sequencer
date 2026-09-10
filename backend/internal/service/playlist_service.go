package service

import (
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/repository"
	ws "github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/websocket"
)

type PlaylistService struct {
	repo *repository.PlaylistRepository
	hub  *ws.Hub
}

func NewPlaylistService(
	repo *repository.PlaylistRepository,
	hub *ws.Hub,
) *PlaylistService {
	return &PlaylistService{
		repo: repo,
		hub:  hub,
	}
}

func (s *PlaylistService) GetByWindowID(
	windowID uint,
) ([]model.PlaylistItem, error) {
	return s.repo.GetByWindowID(windowID)
}

func (s *PlaylistService) Create(
	item *model.PlaylistItem,
) error {
	if err := s.repo.Create(item); err != nil {
		return err
	}

	return BroadcastPlaylistUpdate(s.hub, item.WindowID)
}

func (s *PlaylistService) Delete(
	windowID, itemID uint,
) error {
	if err := s.repo.Delete(windowID, itemID); err != nil {
		return err
	}

	return BroadcastPlaylistUpdate(s.hub, windowID)
}
