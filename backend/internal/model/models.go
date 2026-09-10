package model

import "time"

type Window struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Playlist  []PlaylistItem `json:"playlist,omitempty"`
}

type Media struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"createdAt"`
}

type PlaylistItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	WindowID  uint      `json:"windowId"`
	MediaID   uint      `json:"mediaId"`
	Position  int       `json:"position"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Media Media `gorm:"foreignKey:MediaID" json:"media"`
}

type SyncEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MediaID   uint      `json:"mediaId"`
	StartedAt time.Time `json:"startedAt"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"createdAt"`

	Media Media `json:"media"`
}