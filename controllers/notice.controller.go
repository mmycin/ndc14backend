package controllers

import (
	"net/http"
	"time"

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

func CreateNotice(c *gin.Context) {
	var notice models.Notice

	// Bind JSON input to the notice struct
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	// Set timestamps if your model requires them (optional)
	currentTime := time.Now()
	notice.CreatedAt = currentTime

	// Save the notice to the database
	if err := config.DB.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create notice",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Notice created successfully",
		"data":    notice,
	})
}

func UpdateNotice(c *gin.Context) {
	noticeID := c.Param("id")
	var notice models.Notice

	// Find the existing notice
	if err := config.DB.First(&notice, noticeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notice not found",
		})
		return
	}

	// Bind the updated data
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	// Save the updates
	if err := config.DB.Save(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update notice",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notice updated successfully",
		"data":    notice,
	})
}

func DeleteNotice(c *gin.Context) {
	noticeID := c.Param("id")
	var notice models.Notice

	// Find the notice
	if err := config.DB.First(&notice, noticeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notice not found",
		})
		return
	}

	// Delete the notice and its associated files
	if err := config.DB.Select("Files").Delete(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete notice",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notice deleted successfully",
	})
}

func GetNotice(c *gin.Context) {
	noticeID := c.Param("id")
	var notice models.Notice

	// Find the notice with its associated files
	if err := config.DB.Preload("Files").First(&notice, noticeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notice not found",
		})
		return
	}

	// Censor content of the notice using the generic function
	libs.Censor(&notice)

	c.JSON(http.StatusOK, gin.H{
		"data": notice,
	})
}
