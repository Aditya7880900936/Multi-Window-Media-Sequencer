package repository

import (
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) GetAll() ([]model.Media, error) {
	var media []model.Media

	err := r.db.
		Order("id ASC").
		Find(&media).Error

	return media, err
}

func (r *MediaRepository) GetByID(id uint) (*model.Media, error) {
	var media model.Media

	err := r.db.First(&media, id).Error

	if err != nil {
		return nil, err
	}

	return &media, nil
}

func (r *MediaRepository) Create(media *model.Media) error {
	return r.db.Create(media).Error
}
