package controllers

import (
	"net/http"
	"time"

	"github.com/bluefalconhd/lbd_game/server/database"
	"github.com/bluefalconhd/lbd_game/server/models"
	"github.com/gin-gonic/gin"
)

func GetUserStatistics(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	statistics := make([]map[string]interface{}, 0)
	for _, user := range users {
		var verificationsCount int64
		database.DB.Model(&models.Verification{}).Where("verified_user_id = ?", user.ID).Count(&verificationsCount)

		var phrasesCount int64
		database.DB.Model(&models.Phrase{}).Where("submitted_by = ?", user.ID).Count(&phrasesCount)

		stats := map[string]interface{}{
			"user_id":                user.ID,
			"username":               user.Username,
			"privilege":              user.Privilege,
			"is_eliminated":          user.IsEliminated,
			"verifications_received": verificationsCount,
			"phrases_submitted":      phrasesCount,
		}

		statistics = append(statistics, stats)
	}

	c.JSON(http.StatusOK, gin.H{"statistics": statistics})
}

func ManualReset(c *gin.Context) {
	var input struct {
		OpenTime int64 `json:"open_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetTime := time.Unix(input.OpenTime, 0)

	// Validate that the open time is in the future
	if targetTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Open time must be in the future"})
		return
	}

	// Create new submission window
	window := models.SubmissionWindow{
		OpenTime: targetTime,
	}

	// Start a transaction
	tx := database.DB.Begin()

	// Clean up any existing phrases and verifications for windows that haven't opened yet
	if err := tx.Unscoped().Where("open_time > ?", time.Now()).Delete(&models.SubmissionWindow{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clean up future windows"})
		return
	}

	// Create the new window
	if err := tx.Create(&window).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new submission window"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit changes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Manual reset scheduled",
		"window_id": window.ID,
		"open_time": targetTime,
	})
}

func PromoteUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Privilege = 1 // Promote to Admin Level 1
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted successfully"})
}

func DemoteUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var adminsCount int64
	database.DB.Model(&models.User{}).Where("privilege = ?", 1).Count(&adminsCount)
	if adminsCount == 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot demote the only admin"})
		return
	}

	user.Privilege = 0 // Demote to User
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to demote user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User demoted successfully"})
}

// New endpoint to view scheduled windows
func GetScheduledWindows(c *gin.Context) {
	var windows []models.SubmissionWindow
	if err := database.DB.Where("open_time > ?", time.Now()).Order("open_time asc").Find(&windows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheduled windows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"windows": windows})
}

// New endpoint to cancel a scheduled window
func CancelScheduledWindow(c *gin.Context) {
	windowID := c.Param("id")

	var window models.SubmissionWindow
	if err := database.DB.First(&window, windowID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Window not found"})
		return
	}

	if window.OpenTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel window that has already opened"})
		return
	}

	if err := database.DB.Delete(&window).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel window"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Window cancelled successfully"})
}

func EditPhrase(c *gin.Context) {
	userID := c.GetUint("userID")

	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active submission window"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var phrase models.Phrase
	if err := database.DB.Where("submission_window = ?", window.ID).First(&phrase).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No phrase found for current window"})
		return
	}

	phrase.Content = input.Content
	phrase.SubmittedBy = userID

	if err := database.DB.Save(&phrase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit phrase"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Phrase edited successfully"})
}

func UnsubmitPhrase(c *gin.Context) {
	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active submission window"})
		return
	}

	result := database.DB.Unscoped().Where("submission_window = ?", window.ID).Delete(&models.Phrase{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubmit phrase"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No phrase found for current window"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Phrase unsubmitted successfully"})
}
