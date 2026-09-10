package service

import (
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/repository"
)

type WindowService struct {
	repo *repository.WindowRepository
}

func NewWindowService(repo *repository.WindowRepository) *WindowService {
	return &WindowService{repo: repo}
}

func (s *WindowService) GetAll() ([]model.Window, error) {
	return s.repo.GetAll()
}

func (s *WindowService) GetByID(id uint) (*model.Window, error) {
	return s.repo.GetByID(id)
}