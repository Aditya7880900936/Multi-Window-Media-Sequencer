package repository

import (
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"gorm.io/gorm"
)

type WindowRepository struct {
	db *gorm.DB
}

func NewWindowRepository(db *gorm.DB) *WindowRepository {
	return &WindowRepository{db: db}
}

func (r *WindowRepository) GetAll() ([]model.Window, error) {
	var windows []model.Window

	err := r.db.
		Preload("Playlist.Media").
		Order("id ASC").
		Find(&windows).Error

	return windows, err
}

func (r *WindowRepository) GetByID(id uint) (*model.Window, error) {
	var window model.Window

	err := r.db.
		Preload("Playlist.Media").
		First(&window, id).Error

	if err != nil {
		return nil, err
	}

	return &window, nil
}

func (r *WindowRepository) Create(window *model.Window) error {
	return r.db.Create(window).Error
}