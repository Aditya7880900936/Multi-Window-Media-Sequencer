package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Aditya7880900936/Multi-Window-Media-Sequencer/backend/internal/service"
)

type WindowHandler struct {
	service *service.WindowService
}

func NewWindowHandler(service *service.WindowService) *WindowHandler {
	return &WindowHandler{service: service}
}

func (h *WindowHandler) GetAll(c *gin.Context) {
	windows, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch windows",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": windows,
	})
}

func (h *WindowHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid window id",
		})
		return
	}

	window, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "window not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": window,
	})
}