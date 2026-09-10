package repository

import (
	"time"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"gorm.io/gorm"
)

type PlaybackCycleRepository struct {
	db *gorm.DB
}

func NewPlaybackCycleRepository(db *gorm.DB) *PlaybackCycleRepository {
	return &PlaybackCycleRepository{db: db}
}

func (r *PlaybackCycleRepository) GetOrCreate() (*model.PlaybackCycle, error) {
	var cycle model.PlaybackCycle

	err := r.db.First(&cycle, 1).Error

	if err == nil {
		return &cycle, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	cycle = model.PlaybackCycle{
		ID:        1,
		StartedAt: time.Now().UTC(),
	}

	if err := r.db.Create(&cycle).Error; err != nil {
		return nil, err
	}

	return &cycle, nil
}
