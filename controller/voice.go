package controller

import (
	"log/slog"
	"net/http"

	"go-gin-jwt-demo/config"
	"go-gin-jwt-demo/model"

	"github.com/gin-gonic/gin"
)

type VoiceResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func ListVoices(c *gin.Context) {
	voices := make([]VoiceResponse, 0)

	err := config.DB.WithContext(c.Request.Context()).
		Model(&model.Voice{}).
		Select("id", "name", "gender").
		Where("is_available = ?", true).
		Order("id ASC").
		Find(&voices).Error

	if err != nil {
		slog.Error("voice_list_failed", "error", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Unable to load voices. Please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"voices": voices,
	})
}
