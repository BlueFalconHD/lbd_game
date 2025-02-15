package controllers

import (
	"net/http"
	"time"

	"github.com/bluefalconhd/lbd_game/server/database"
	"github.com/bluefalconhd/lbd_game/server/models"
	"github.com/gin-gonic/gin"
)

func VerifyUser(c *gin.Context) {
	verifierID := c.GetUint("userID")

	// Get current submission window
	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active submission window"})
		return
	}

	var input struct {
		VerifiedUserID uint `json:"verified_user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if verification already exists for this window
	var existingVerification models.Verification
	if err := database.DB.Where("verified_user_id = ? AND submission_window = ?",
		input.VerifiedUserID, window.ID).First(&existingVerification).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User has already been verified in this window"})
		return
	}

	verification := models.Verification{
		VerifiedUserID:   input.VerifiedUserID,
		VerifierID:       verifierID,
		SubmissionWindow: window.ID,
	}

	if err := database.DB.Create(&verification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record verification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Verification recorded"})
}

func GetCurrentVerifications(c *gin.Context) {
	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "No active submission window"})
		return
	}

	var verifications []struct {
		VerificationID   uint      `json:"verification_id"`
		VerifierID       uint      `json:"verifier_id"`
		VerifierName     string    `json:"verifier_name"`
		VerifiedID       uint      `json:"verified_id"`
		VerifiedName     string    `json:"verified_name"`
		SubmissionWindow uint      `json:"submission_window"`
		CreatedAt        time.Time `json:"created_at"`
	}

	result := database.DB.Table("verifications").
		Select("verifications.id as verification_id, verifications.verifier_id, "+
			"u1.username as verifier_name, verifications.verified_user_id as verified_id, "+
			"u2.username as verified_name, verifications.submission_window, verifications.created_at").
		Joins("JOIN users u1 ON verifications.verifier_id = u1.id").
		Joins("JOIN users u2 ON verifications.verified_user_id = u2.id").
		Where("verifications.submission_window = ?", window.ID).
		Find(&verifications)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch verifications"})
		return
	}

	c.JSON(http.StatusOK, verifications)
}

func GetUnverifiedUsers(c *gin.Context) {
	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "No active submission window"})
		return
	}

	var unverifiedUsers []struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	result := database.DB.Table("users").
		Select("id, username").
		Where("id NOT IN (?)",
			database.DB.Table("verifications").
				Select("verified_user_id").
				Where("submission_window = ?", window.ID)).
		Find(&unverifiedUsers)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch unverified users"})
		return
	}

	c.JSON(http.StatusOK, unverifiedUsers)
}
