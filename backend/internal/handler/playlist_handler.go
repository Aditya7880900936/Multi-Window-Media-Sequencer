package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/model"
	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/service"
)

type PlaylistHandler struct {
	service *service.PlaylistService
}

func NewPlaylistHandler(service *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service: service}
}

func (h *PlaylistHandler) GetByWindowID(c *gin.Context) {
	windowID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid window id",
		})
		return
	}

	playlist, err := h.service.GetByWindowID(uint(windowID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch playlist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": playlist,
	})
}

func (h *PlaylistHandler) Create(c *gin.Context) {
	windowID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid window id",
		})
		return
	}

	var item model.PlaylistItem

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	item.WindowID = uint(windowID)

	if err := h.service.Create(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add playlist item",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": item,
	})
}

func (h *PlaylistHandler) Delete(c *gin.Context) {
	windowID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid window id",
		})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid playlist item id",
		})
		return
	}

	if err := h.service.Delete(uint(windowID), uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete playlist item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "playlist item deleted",
	})
}