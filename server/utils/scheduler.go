package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bluefalconhd/lbd_game/server/database"
	"github.com/bluefalconhd/lbd_game/server/models"
	"github.com/robfig/cron/v3"
)

var (
	cst *time.Location
)

func init() {
	var err error
	cst, err = time.LoadLocation("America/Chicago")
	if err != nil {
		panic(err)
	}
}

func InitScheduler() {
	c := cron.New(cron.WithLocation(cst))
	// Schedule at midnight CST
	c.AddFunc("0 0 * * *", scheduleSubmissionWindow)

	// Run immediately if there's no window for today
	today := time.Now().In(cst).Truncate(24 * time.Hour)
	var window models.SubmissionWindow
	if database.DB.Where("date = ?", today).First(&window).Error != nil {
		scheduleSubmissionWindow()
	}

	c.Start()
}

func IsSubmissionWindowOpen() bool {
	now := time.Now().In(cst)

	// Get most recent submission window
	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		return false
	}

	// Check if we've passed the open time
	if now.Before(window.OpenTime) {
		return false
	}

	// Check if a phrase has already been submitted for this window
	var phrase models.Phrase
	if err := database.DB.Where("submission_window = ?", window.ID).First(&phrase).Error; err == nil {
		return false
	}

	return true
}

func randomTime(start, end time.Time) time.Time {
	delta := end.Sub(start)
	sec := rand.Int63n(int64(delta.Seconds()))
	return start.Add(time.Duration(sec) * time.Second)
}

func GetCurrentWindow() (*models.SubmissionWindow, error) {
	var window models.SubmissionWindow
	if err := database.DB.Order("open_time desc").First(&window).Error; err != nil {
		return nil, err
	}
	return &window, nil
}

func GetNextScheduledWindow() (*models.SubmissionWindow, error) {
	var window models.SubmissionWindow
	if err := database.DB.Where("open_time > ?", time.Now()).
		Order("open_time asc").
		First(&window).Error; err != nil {
		return nil, err
	}
	return &window, nil
}

func CleanupOldWindows() error {
	// Keep only the last 30 days of windows
	cutoff := time.Now().AddDate(0, 0, -30)
	return database.DB.Unscoped().
		Where("open_time < ?", cutoff).
		Delete(&models.SubmissionWindow{}).Error
}

// Modify the existing scheduleSubmissionWindow function
func scheduleSubmissionWindow() {
	// Check if there's already a window scheduled
	if window, err := GetNextScheduledWindow(); err == nil && window != nil {
		// Already have a scheduled window
		return
	}

	now := time.Now().In(cst)
	tomorrow := now.Add(24 * time.Hour).Truncate(24 * time.Hour)

	// Generate random time between 4:30 AM and 8:20 AM CST for tomorrow
	startRange := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 4, 30, 0, 0, cst)
	endRange := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 8, 20, 0, 0, cst)
	openTime := randomTime(startRange, endRange)

	window := models.SubmissionWindow{
		OpenTime: openTime,
	}

	if err := database.DB.Create(&window).Error; err != nil {
		// Log the error appropriately
		fmt.Printf("Failed to schedule submission window: %v\n", err)
	}

	// Cleanup old windows
	CleanupOldWindows()
}
