package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/config"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/database"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
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

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
