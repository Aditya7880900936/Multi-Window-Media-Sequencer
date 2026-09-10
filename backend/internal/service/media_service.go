package service

import (
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/repository"
)

type MediaService struct {
	repo *repository.MediaRepository
}

func NewMediaService(repo *repository.MediaRepository) *MediaService {
	return &MediaService{repo: repo}
}

func (s *MediaService) GetAll() ([]model.Media, error) {
	return s.repo.GetAll()
}

func (s *MediaService) GetByID(id uint) (*model.Media, error) {
	return s.repo.GetByID(id)
}

func (s *MediaService) Create(media *model.Media) error {
	return s.repo.Create(media)
}