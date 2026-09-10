package repository

import (
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"gorm.io/gorm"
)

type PlaylistRepository struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) *PlaylistRepository {
	return &PlaylistRepository{db: db}
}

func (r *PlaylistRepository) GetByWindowID(windowID uint) ([]model.PlaylistItem, error) {
	var items []model.PlaylistItem

	err := r.db.
		Where("window_id = ?", windowID).
		Preload("Media").
		Order("position ASC").
		Find(&items).Error

	return items, err
}

func (r *PlaylistRepository) Create(item *model.PlaylistItem) error {
	return r.db.Create(item).Error
}

func (r *PlaylistRepository) Delete(windowID, itemID uint) error {
	return r.db.
		Where("window_id = ? AND id = ?", windowID, itemID).
		Delete(&model.PlaylistItem{}).Error
}