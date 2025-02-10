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
	"github.com/gin-contrib/cors"

	// to read from .env file
	"github.com/joho/godotenv"
)

var (
    db                     *gorm.DB
    jwtKey                 = []byte("")
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
    IsAdmin	    bool   `gorm:"default:false"`
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

    // Load environment variables
    err := godotenv.Load()
    if err != nil {
    	log.Println("Failed to load .env file")
     	panic(err)
    }

    jwtKey = []byte(os.Getenv("JWT_SECRET"))
    frontendURL := os.Getenv("FRONTEND_URL")

    if jwtKey == nil {
    	log.Fatal("JWT_SECRET environment variable not set")
    }

    if frontendURL == "" {
    	log.Fatal("FRONTEND_URL environment variable not set")
    }

    // Initialize the database and schedule submission time
    initDatabase()
    schedulePhraseSubmission()

    // Initialize Gin router
    router := gin.Default()

    // CORS configuration
    config := cors.Config{
        AllowOrigins:     []string{frontendURL},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }

    router.Use(cors.New(config))

    // Public routes
    router.POST("/register", registerUser)
    router.POST("/login", loginUser)

    userRoutes := router.Group("/user")
    {
        userRoutes.Use(authMiddleware(false))
        userRoutes.GET("/status", getUserStatus)
        userRoutes.POST("/submit_phrase", submitPhrase)
        userRoutes.POST("/confirm_usage", confirmPhraseUsage)
        userRoutes.GET("/phrase", getCurrentPhrase)
        userRoutes.GET("/verification_status", getVerificationStatus)
        userRoutes.GET("/active_users", getActiveUsers)
        userRoutes.GET("/verifications", getTodaysVerifications)
    }

    // Additional route for admin
    adminRoutes := router.Group("/admin")
    {
        adminRoutes.Use(authMiddleware(true))

        adminRoutes.POST("/register_admin", registerAdmin)
        adminRoutes.GET("/pending_users", getPendingUsers)

        adminRoutes.POST("/edit_phrase", editPhrase)
        adminRoutes.POST("/reset_game", resetGame)

        adminRoutes.POST("/approve_user", approveUser)
        adminRoutes.POST("/eliminate_user", eliminateUser)
        adminRoutes.POST("/resurrect_user", resurrectUser)
        adminRoutes.POST("/unapprove_user", unapproveUser)
        adminRoutes.POST("/set_admin", setAdmin)
        adminRoutes.GET("/detailed_users", getDetailedUsers)
    }

    // Start the server
    router.Run(":8040")
}

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
    db.AutoMigrate(&User{}, &Phrase{}, &PhraseUsage{})

    // Create default admin if not exists
    // var admin User
    // if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
    //     passwordHash, _ := bcrypt.GenerateFromPassword([]byte("adminpass"), bcrypt.DefaultCost)
    //     admin = User{
    //         Username:   "admin",
    //         Password:   string(passwordHash),
    //         IsApproved: true,
    //         IsAdmin:    true,
    //     }
    //     db.Create(&admin)
    //     log.Println("Default admin created with username 'admin' and password 'adminpass'")
    // }
}

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

        var user User
        if err := db.First(&user, claims.UserID).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if requireAdmin && !user.IsAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
            c.Abort()
            return
        }

        // Update claims with latest IsAdmin status
        claims.IsAdmin = user.IsAdmin

        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("isAdmin", claims.IsAdmin)

        c.Next()
    }
}

func getPendingUsers(c *gin.Context) {
    var users []User
    db.Where("is_approved = ?", false).Find(&users)

    var userList []gin.H
    for _, user := range users {
        userList = append(userList, gin.H{"username": user.Username})
    }

    c.JSON(http.StatusOK, gin.H{"users": userList})
}

func getTodaysVerifications(c *gin.Context) {
    today := time.Now().Truncate(24 * time.Hour)

    var usages []PhraseUsage
    db.Where("date = ?", today).Find(&usages)

    var verifications []gin.H
    for _, usage := range usages {
        var user, confirmer User
        db.First(&user, usage.UserID)
        db.First(&confirmer, usage.ConfirmedByUserID)

        verifications = append(verifications, gin.H{
            "userId":      user.ID,
            "username":    user.Username,
            "confirmedBy": confirmer.Username,
        })
    }

    c.JSON(http.StatusOK, gin.H{"verifications": verifications})
}

func getVerificationStatus(c *gin.Context) {
    userIDRaw, _ := c.Get("userID")
    userID := userIDRaw.(uint)
    today := time.Now().Truncate(24 * time.Hour)

    var usage PhraseUsage
    result := db.Where("user_id = ? AND date = ?", userID, today).First(&usage)
    if result.Error != nil {
        c.JSON(http.StatusOK, gin.H{
            "verified":   false,
            "verified_by": "",
        })
        return
    }

    var verifier User
    db.First(&verifier, usage.ConfirmedByUserID)

    c.JSON(http.StatusOK, gin.H{
        "verified":   true,
        "verified_by": verifier.Username,
    })
}

func getActiveUsers(c *gin.Context) {
    userIDRaw, _ := c.Get("userID")
    currentUserID := userIDRaw.(uint)
    var users []User
    db.Where("is_eliminated = ? AND is_approved = ?", false, true).Find(&users)

    var userList []gin.H
    for _, user := range users {
        if user.ID != currentUserID {
            userList = append(userList, gin.H{"username": user.Username})
        }
    }

    c.JSON(http.StatusOK, gin.H{"users": userList})
}

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
        IsAdmin:    false, // Regular user
    }
    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
        return
    }

    log.Printf("New user registered: %s", user.Username)
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully. Awaiting admin approval."})
}

func loginUser(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if user exists
    var user User
    if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // Compare passwords
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    if !user.IsApproved && !user.IsAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "Account not approved by admin"})
        return
    }

    if user.IsEliminated {
        c.JSON(http.StatusForbidden, gin.H{"error": "You have been eliminated"})
        return
    }

    // Generate token
    expirationTime := time.Now().Add(24 * time.Hour)

    claims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        IsAdmin:  user.IsAdmin,
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

    log.Printf("User '%s' logged in", user.Username)
    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

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

    admin := User{
    	Username:   input.Username,
     	Password:   string(passwordHash),
       	IsApproved: true,
        IsAdmin:    true,
    }
    if err := db.Create(&admin).Error; err != nil {
    	c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
        return
    }

    log.Printf("New admin registered: %s", admin.Username)
    c.JSON(http.StatusCreated, gin.H{"message": "Admin registered successfully"})
}

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

func getCurrentPhrase(c *gin.Context) {
    today := time.Now().Truncate(24 * time.Hour)
    var phrase Phrase
    result := db.Where("date = ?", today).First(&phrase)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No phrase submitted today"})
        return
    }

    var submittedBy User
    db.First(&submittedBy, phrase.SubmittedByID)

    c.JSON(http.StatusOK, gin.H{
        "phrase": gin.H{
            "text":        phrase.Text,
            "submittedBy": submittedBy.Username,
        },
    })
}

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

func editPhrase(c *gin.Context) {
	var input struct {
		Text string `json:"text" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	var phrase Phrase
	result := db.Where("date = ?", today).First(&phrase)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No phrase submitted today"})
		return
	}

	phrase.Text = input.Text
	db.Save(&phrase)

	log.Printf("Phrase edited by admin: %s", phrase.Text)
	c.JSON(http.StatusOK, gin.H{"message": "Phrase edited successfully"})
}

func eliminateUser(c *gin.Context) {
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

	user.IsEliminated = true
	db.Save(&user)

	log.Printf("User '%s' eliminated by admin", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "User eliminated successfully"})
}

func resetGame(c *gin.Context) {
	dailyReset()
	c.JSON(http.StatusOK, gin.H{"message": "Game reset successfully"})
}

func unapproveUser(c *gin.Context) {
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

	user.IsApproved = false
	db.Save(&user)

	log.Printf("User '%s' unapproved by admin", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "User unapproved successfully"})
}

func resurrectUser(c *gin.Context) {
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

	user.IsEliminated = false
	db.Save(&user)

	log.Printf("User '%s' resurrected by admin", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "User resurrected successfully"})
}

func setAdmin(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Admin   bool   `json:"admin" binding:"required"`
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

	user.IsAdmin = input.Admin
	db.Save(&user)

	log.Printf("User '%s' promoted to admin by admin", user.Username)
	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}

func getDetailedUsers(c *gin.Context) {
	var users []User
	db.Find(&users)

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"is_approved": user.IsApproved,
			"is_eliminated": user.IsEliminated,
			"is_admin":    user.IsAdmin,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": userList})
}
