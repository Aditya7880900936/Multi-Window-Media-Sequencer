package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/service"
)

type SyncHandler struct {
	service *service.SyncService
}

func NewSyncHandler(service *service.SyncService) *SyncHandler {
	return &SyncHandler{
		service: service,
	}
}

type SyncRequest struct {
	MediaID uint `json:"mediaId" binding:"required"`
	Duration int `json:"duration" binding:"required"`
}

func (h *SyncHandler) Sync(c *gin.Context) {
	var req SyncRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "mediaId and duration are required",
		})
		return
	}

	if err := h.service.Sync(req.MediaID, req.Duration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sync media",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sync started",
	})
}