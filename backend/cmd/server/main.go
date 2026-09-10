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
	)
	if err != nil {
		panic(err)
	}

	if err := database.Seed(db); err != nil {
		panic(err)
	}

	// Repositories
	windowRepo := repository.NewWindowRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	// Services
	windowService := service.NewWindowService(windowRepo)
	mediaService := service.NewMediaService(mediaRepo)
	playlistService := service.NewPlaylistService(playlistRepo)

	// Handlers
	windowHandler := handler.NewWindowHandler(windowService)
	mediaHandler := handler.NewMediaHandler(mediaService)
	playlistHandler := handler.NewPlaylistHandler(playlistService)

	router := gin.Default()

	// Health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API
	api := router.Group("/api")

	// Windows
	api.GET("/windows", windowHandler.GetAll)
	api.GET("/windows/:id", windowHandler.GetByID)

	// Media
	api.GET("/media", mediaHandler.GetAll)
	api.POST("/media", mediaHandler.Create)

	// Playlists
	api.GET("/windows/:id/playlist", playlistHandler.GetByWindowID)
	api.POST("/windows/:id/playlist", playlistHandler.Create)
	api.DELETE("/windows/:id/playlist/:itemId", playlistHandler.Delete)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}