package repository

import (
	"time"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ActiveSyncRepository struct {
	db *gorm.DB
}

func NewActiveSyncRepository(db *gorm.DB) *ActiveSyncRepository {
	return &ActiveSyncRepository{db: db}
}

func (r *ActiveSyncRepository) Get() (*model.ActiveSync, error) {
	var sync model.ActiveSync

	err := r.db.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Silent),
	}).Preload("Media").First(&sync, 1).Error

	if err != nil {
		return nil, err
	}

	return &sync, nil
}

func (r *ActiveSyncRepository) CreateOrUpdate(
	mediaID uint,
	startedAt time.Time,
	duration int,
) error {
	var sync model.ActiveSync

	err := r.db.First(&sync, 1).Error

	if err == gorm.ErrRecordNotFound {
		sync = model.ActiveSync{
			ID:        1,
			MediaID:   mediaID,
			StartedAt: startedAt,
			Duration:  duration,
		}

		return r.db.Create(&sync).Error
	}

	if err != nil {
		return err
	}

	sync.MediaID = mediaID
	sync.StartedAt = startedAt
	sync.Duration = duration

	return r.db.Save(&sync).Error
}

func (r *ActiveSyncRepository) Clear() error {
	return r.db.Delete(&model.ActiveSync{}, 1).Error
}
