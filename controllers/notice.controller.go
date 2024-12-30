package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmycin/ndc14/config"
	"github.com/mmycin/ndc14/libs"
	"github.com/mmycin/ndc14/models"
)

func GetNotices(c *gin.Context) {
	var notices []models.Notice

	// Get all notices with their associated files, ordered by creation date (newest first)
	if err := config.DB.Preload("Files").Order("created_at DESC").Find(&notices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Unable to fetch notices",
			"details": err.Error(),
		})
		return
	}

	// Censor content of notices using the generic function
	libs.Censor(&notices)

	libs.ReverseArray(notices)

	c.JSON(http.StatusOK, gin.H{
		"data":  notices,
		"count": len(notices),
	})
}
