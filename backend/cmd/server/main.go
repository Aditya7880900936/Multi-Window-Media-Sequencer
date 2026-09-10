package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/config"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/database"
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