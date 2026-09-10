package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/service"
)

type MediaHandler struct {
	service *service.MediaService
}

func NewMediaHandler(service *service.MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

func (h *MediaHandler) GetAll(c *gin.Context) {
	media, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch media",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": media,
	})
}

func (h *MediaHandler) Create(c *gin.Context) {
	var media model.Media

	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.service.Create(&media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create media",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": media,
	})
}
