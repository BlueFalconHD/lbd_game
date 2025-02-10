package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
    db                     *gorm.DB
    jwtKey                 = []byte("abcdefg") // Change this to a secure key
    phraseSubmissionOpen   bool
    phraseSubmissionTime   time.Time
    currentPhrase          Phrase
    dailyResetScheduled    bool
    totalUsers             int
    logFile                *os.File
)

type User struct {
    ID          uint   `gorm:"primaryKey"`
    Username    string `gorm:"unique;not null"`
    Password    string `gorm:"not null"`
    IsApproved  bool   `gorm:"default:false"`
    IsEliminated bool   `gorm:"default:false"`
}

type Admin struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
}

type Phrase struct {
    ID            uint      `gorm:"primaryKey"`
    Date          time.Time `gorm:"unique;not null"`
    Text          string    `gorm:"not null"`
    SubmittedByID uint      `gorm:"not null"`
    SubmissionTime time.Time
}

type PhraseUsage struct {
    ID              uint      `gorm:"primaryKey"`
    UserID          uint      `gorm:"not null"`
    ConfirmedByUserID uint    `gorm:"not null"`
    Date            time.Time `gorm:"not null"`
}

type Claims struct {
    UserID   uint `json:"user_id"`
    Username string `json:"username"`
    IsAdmin  bool `json:"is_admin"`
    jwt.StandardClaims
}

func main() {
    initLogging()
    defer logFile.Close()


    // Initialize the database and schedule submission time
    initDatabase()
    schedulePhraseSubmission()

    // Initialize Gin router
    router := gin.Default()

    router.LoadHTMLGlob("templates/*")
    router.Static("/static", "./static")

    // Public routes
    router.POST("/register", registerUser)
    router.POST("/login", loginUser)

    // Admin routes
    adminRoutes := router.Group("/admin")
    {
        adminRoutes.Use(authMiddleware(true))

        adminRoutes.POST("/approve_user", approveUser)
        adminRoutes.POST("/register_admin", registerAdmin)
    }

    // Protected user routes
    userRoutes := router.Group("/user")
    {
        userRoutes.Use(authMiddleware(false))
        userRoutes.GET("/status", getUserStatus)
        userRoutes.POST("/submit_phrase", submitPhrase)
        userRoutes.POST("/confirm_usage", confirmPhraseUsage)
    }

    // Start the server
    router.Run(":8080")
}

// Initialize Logging
func initLogging() {
    var err error
    logFile, err = os.OpenFile("lbd_game.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("Failed to open log file:", err)
    }
    log.SetOutput(logFile)
}

func initDatabase() {
    var err error
    db, err = gorm.Open(sqlite.Open("lbd_game.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect database: %v", err)
    }

    // Migrate schema
    db.AutoMigrate(&User{}, &Admin{}, &Phrase{}, &PhraseUsage{})

    // Create default admin if not exists
    var adminCount int64
    db.Model(&Admin{}).Count(&adminCount)

    log.Println("Admin count:", adminCount)
    if adminCount == 0 {
        passwordHash, _ := bcrypt.GenerateFromPassword([]byte("adminpass"), bcrypt.DefaultCost)
        admin := Admin{
            Username: "admin",
            Password: string(passwordHash),
        }
        db.Create(&admin)
        log.Println("Default admin created with username 'admin' and password 'adminpass'")
    }
}

// Schedule Phrase Submission
func schedulePhraseSubmission() {
    now := time.Now()
    year, month, day := now.Date()
    location := now.Location()

    // Define the time window
    startTime := time.Date(year, month, day, 4, 30, 0, 0, location)
    endTime := time.Date(year, month, day, 8, 20, 0, 0, location)

    // Random time within the window
    rand.Seed(time.Now().UnixNano())
    randomDuration := time.Duration(rand.Intn(int(endTime.Sub(startTime).Seconds()))) * time.Second
    phraseSubmissionTime = startTime.Add(randomDuration)

    // Schedule the opening of the submission window
    timeUntilOpen := time.Until(phraseSubmissionTime)
    log.Printf("Phrase submission will open at %s", phraseSubmissionTime.Format(time.Kitchen))
    time.AfterFunc(timeUntilOpen, func() {
        phraseSubmissionOpen = true
        log.Println("Phrase submission is now open!")
    })

    // Schedule daily reset at midnight if not already scheduled
    if !dailyResetScheduled {
        midnight := time.Date(year, month, day+1, 0, 0, 0, 0, location)
        timeUntilMidnight := time.Until(midnight)
        time.AfterFunc(timeUntilMidnight, dailyReset)
        dailyResetScheduled = true
        log.Println("Daily reset scheduled at midnight")
    }
}

// Daily Reset Function
func dailyReset() {
    log.Println("Performing daily reset")

    // Elimination logic
    today := time.Now().Truncate(24 * time.Hour)
    var users []User
    db.Where("is_eliminated = ? AND is_approved = ?", false, true).Find(&users)
    var activeUsers int
    for _, user := range users {
        // Check if user had usage confirmed today
        var usage PhraseUsage
        result := db.Where("user_id = ? AND date = ?", user.ID, today).First(&usage)
        if result.Error != nil {
            // User did not have usage confirmed; eliminate
            user.IsEliminated = true
            db.Save(&user)
            log.Printf("User '%s' has been eliminated", user.Username)
        } else {
            activeUsers++
        }
    }

    // Check for winners
    if activeUsers <= 3 {
        var winners []User
        db.Where("is_eliminated = ? AND is_approved = ?", false, true).Find(&winners)
        for _, winner := range winners {
            log.Printf("User '%s' is a winner!", winner.Username)
        }
        // Reset the game for next time (optional)
    }

    // Clear today's PhraseUsage records
    db.Where("date = ?", today).Delete(&PhraseUsage{})

    // Reset phrase submission
    phraseSubmissionOpen = false
    schedulePhraseSubmission()
    dailyResetScheduled = false
    log.Println("Daily reset complete")
}

// Authentication Middleware
func authMiddleware(requireAdmin bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            c.Abort()
            return
        }

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if requireAdmin && !claims.IsAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("isAdmin", claims.IsAdmin)

        c.Next()
    }
}

// Register User
func registerUser(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    user := User{
        Username:   input.Username,
        Password:   string(passwordHash),
        IsApproved: false, // Needs admin approval
    }
    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
        return
    }

    log.Printf("New user registered: %s", user.Username)
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully. Awaiting admin approval."})
}

// Login User
// Login User or Admin
func loginUser(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var isAdmin bool
    var userID uint
    var username string

    // Attempt to find user in the Users table
    var user User
    if err := db.Where("username = ?", input.Username).First(&user).Error; err == nil {
        // Compare passwords
        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }

        if !user.IsApproved {
            c.JSON(http.StatusForbidden, gin.H{"error": "Account not approved by admin"})
            return
        }

        if user.IsEliminated {
            c.JSON(http.StatusForbidden, gin.H{"error": "You have been eliminated"})
            return
        }

        isAdmin = false
        userID = user.ID
        username = user.Username

    } else {
        // If not found in Users table, attempt to find in Admins table
        var admin Admin
        if err := db.Where("username = ?", input.Username).First(&admin).Error; err == nil {
            // Compare passwords
            if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
                return
            }

            isAdmin = true
            userID = admin.ID
            username = admin.Username

        } else {
            // User not found in either table
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }
    }

    // Generate token
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID:   userID,
        Username: username,
        IsAdmin:  isAdmin,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
        return
    }

    if isAdmin {
        log.Printf("Admin '%s' logged in", username)
    } else {
        log.Printf("User '%s' logged in", username)
    }
    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Register Admin
func registerAdmin(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    admin := Admin{
        Username: input.Username,
        Password: string(passwordHash),
    }
    if err := db.Create(&admin).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Admin username already exists"})
        return
    }

    log.Printf("New admin registered: %s", admin.Username)
    c.JSON(http.StatusCreated, gin.H{"message": "Admin registered successfully"})
}

// Approve User
func approveUser(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user User
    if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    user.IsApproved = true
    db.Save(&user)

    log.Printf("User '%s' approved by admin", user.Username)
    c.JSON(http.StatusOK, gin.H{"message": "User approved successfully"})
}

// Get User Status
func getUserStatus(c *gin.Context) {
    userIDRaw, _ := c.Get("userID")
    userID := userIDRaw.(uint)

    var user User
    if err := db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Check if user is eliminated
    isEliminated := user.IsEliminated

    c.JSON(http.StatusOK, gin.H{
        "user_id":       user.ID,
        "username":      user.Username,
        "is_eliminated": isEliminated,
    })
}

// Submit Phrase
func submitPhrase(c *gin.Context) {
    if !phraseSubmissionOpen {
        c.JSON(http.StatusForbidden, gin.H{"error": "Phrase submission is not open"})
        return
    }

    userIDRaw, _ := c.Get("userID")
    userID := userIDRaw.(uint)

    var input struct {
        Text string `json:"text" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if a phrase has already been submitted today
    var existingPhrase Phrase
    today := time.Now().Truncate(24 * time.Hour)
    result := db.Where("date = ?", today).First(&existingPhrase)
    if result.Error == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Phrase already submitted for today"})
        return
    }

    phrase := Phrase{
        Date:          today,
        Text:          input.Text,
        SubmittedByID: userID,
        SubmissionTime: time.Now(),
    }
    db.Create(&phrase)
    currentPhrase = phrase

    // Close submission window
    phraseSubmissionOpen = false

    log.Printf("Phrase submitted by user '%s': %s", c.GetString("username"), phrase.Text)
    c.JSON(http.StatusOK, gin.H{"message": "Phrase submitted successfully", "phrase": phrase.Text})
}

// Confirm Phrase Usage
func confirmPhraseUsage(c *gin.Context) {
    userIDRaw, _ := c.Get("userID")
    confirmerID := userIDRaw.(uint)

    var input struct {
        Username string `json:"username" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Find the user to confirm
    var user User
    if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Cannot confirm self
    if user.ID == confirmerID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot confirm yourself"})
        return
    }

    // Check if already confirmed today
    today := time.Now().Truncate(24 * time.Hour)
    var existingUsage PhraseUsage
    result := db.Where("user_id = ? AND date = ?", user.ID, today).First(&existingUsage)
    if result.Error == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Phrase usage already confirmed for this user today"})
        return
    }

    usage := PhraseUsage{
        UserID:          user.ID,
        ConfirmedByUserID: confirmerID,
        Date:            today,
    }
    db.Create(&usage)

    log.Printf("User '%s' confirmed phrase usage for '%s'", c.GetString("username"), user.Username)
    c.JSON(http.StatusOK, gin.H{"message": "Phrase usage confirmed"})
}
