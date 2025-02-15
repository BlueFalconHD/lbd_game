package controllers

import (
	"net/http"
	"time"

	"github.com/bluefalconhd/lbd_game/server/database"
	"github.com/bluefalconhd/lbd_game/server/models"
	"github.com/bluefalconhd/lbd_game/server/utils"
	"github.com/gin-gonic/gin"
)

func GetCurrentPhrase(c *gin.Context) {
	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"phrase": nil, "message": "No active submission window"})
		return
	}

	var phrase models.Phrase
	if err := database.DB.Where("submission_window = ?", window.ID).First(&phrase).Error; err != nil {
		nextOpenTime := window.OpenTime
		if time.Now().After(window.OpenTime) {
			// If we're past the open time and no phrase exists, show next window
			nextWindow := models.SubmissionWindow{}
			if err := database.DB.Where("open_time > ?", window.OpenTime).Order("open_time").First(&nextWindow).Error; err == nil {
				nextOpenTime = nextWindow.OpenTime
			}
		}
		c.JSON(http.StatusOK, gin.H{
			// "phrase":         nil,
			"message":        "No phrase submitted yet",
			"next_open_time": nextOpenTime,
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, phrase.SubmittedBy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"phrase":            phrase.Content,
		"submittedBy":       user.Username,
		"submission_window": phrase.SubmissionWindow,
	})
}

func SubmitPhrase(c *gin.Context) {
	userID := c.GetUint("userID")

	if !utils.IsSubmissionWindowOpen() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Submission window is closed"})
		return
	}

	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No active submission window"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phrase := models.Phrase{
		Content:          input.Content,
		SubmittedBy:      userID,
		SubmissionWindow: window.ID,
	}

	if err := database.DB.Create(&phrase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit phrase"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Phrase submitted successfully"})
}

func CanSubmitPhrase(c *gin.Context) {
	if utils.IsSubmissionWindowOpen() {
		c.JSON(http.StatusOK, gin.H{"can_submit": true})
		return
	}

	c.JSON(http.StatusOK, gin.H{"can_submit": false})
}
