package database

import (
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64

	if err := db.Model(&model.Window{}).Count(&count).Error; err != nil {
		return err
	}

	// Don't seed again if data already exists.
	if count > 0 {
		return nil
	}

	// Media
	media := []model.Media{
		{
			Name:     "Image 1",
			Type:     "image",
			URL:      "https://picsum.photos/id/10/1280/720",
			Duration: 30,
		},
		{
			Name:     "Image 2",
			Type:     "image",
			URL:      "https://picsum.photos/id/20/1280/720",
			Duration: 20,
		},
		{
			Name:     "Video 1",
			Type:     "video",
			URL:      "https://www.w3schools.com/html/mov_bbb.mp4",
			Duration: 40,
		},
	}

	if err := db.Create(&media).Error; err != nil {
		return err
	}

	// Windows
	windows := []model.Window{
		{Name: "Window 1"},
		{Name: "Window 2"},
		{Name: "Window 3"},
	}

	if err := db.Create(&windows).Error; err != nil {
		return err
	}

	// Window 1 playlist
	playlist := []model.PlaylistItem{
		{
			WindowID: windows[0].ID,
			MediaID:  media[0].ID,
			Position: 1,
			Duration: media[0].Duration,
		},
		{
			WindowID: windows[0].ID,
			MediaID:  media[1].ID,
			Position: 2,
			Duration: media[1].Duration,
		},
		{
			WindowID: windows[0].ID,
			MediaID:  media[2].ID,
			Position: 3,
			Duration: media[2].Duration,
		},

		// Window 2
		{
			WindowID: windows[1].ID,
			MediaID: media[1].ID,
			Position: 1,
			Duration: media[1].Duration,
		},
		{
			WindowID: windows[1].ID,
			MediaID: media[2].ID,
			Position: 2,
			Duration: media[2].Duration,
		},

		// Window 3
		{
			WindowID: windows[2].ID,
			MediaID: media[2].ID,
			Position: 1,
			Duration: media[2].Duration,
		},
		{
			WindowID: windows[2].ID,
			MediaID: media[0].ID,
			Position: 2,
			Duration: media[0].Duration,
		},
	}

	return db.Create(&playlist).Error
}