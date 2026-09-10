package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/config"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/database"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/handler"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/repository"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/service"
	ws "github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/websocket"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	err = db.AutoMigrate(
		&model.Window{},
		&model.Media{},
		&model.PlaylistItem{},
		&model.SyncEvent{},
		&model.PlaybackCycle{},
		&model.ActiveSync{},
	)
	if err != nil {
		panic(err)
	}

	if err := database.Seed(db); err != nil {
		panic(err)
	}

	// WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Repositories
	windowRepo := repository.NewWindowRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)
	cycleRepo := repository.NewPlaybackCycleRepository(db)

	// Services
	sequencer := service.NewSequencer()

	// Sync Service
	syncService := service.NewSyncService(db, hub)

	windowService := service.NewWindowService(
		windowRepo,
		sequencer,
		cycleRepo,
		syncService,
	)
	mediaService := service.NewMediaService(mediaRepo)
	playlistService := service.NewPlaylistService(
		playlistRepo,
		hub,
	)

	// Handlers
	windowHandler := handler.NewWindowHandler(windowService)
	mediaHandler := handler.NewMediaHandler(mediaService)
	playlistHandler := handler.NewPlaylistHandler(playlistService)
	websocketHandler := handler.NewWebSocketHandler(hub)
	syncHandler := handler.NewSyncHandler(syncService)

	router := gin.Default()

	// Health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// WebSocket
	router.GET("/ws", websocketHandler.Handle)

	// API
	api := router.Group("/api")

	// Windows
	api.GET("/windows", windowHandler.GetAll)
	api.GET("/windows/:id", windowHandler.GetByID)
	api.GET("/windows/:id/current", windowHandler.GetCurrentPlayback)

	// Media
	api.GET("/media", mediaHandler.GetAll)
	api.POST("/media", mediaHandler.Create)

	// Playlists
	api.GET("/windows/:id/playlist", playlistHandler.GetByWindowID)
	api.POST("/windows/:id/playlist", playlistHandler.Create)
	api.DELETE("/windows/:id/playlist/:itemId", playlistHandler.Delete)

	// Sync
	api.POST("/sync", syncHandler.Sync)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
