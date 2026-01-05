package handlers

import (
	"goback/config"
	"goback/database"
	"goback/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	input.Password = string(hashedPassword)

	// Create the User
	if result := database.DB.Create(&input); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email likely already exists"})
		return
	}

	defaultTasks := []models.Task{
		{
			Title:    "Personal tasks and notes",
			Category: "General",
			Priority: 3,
			UserID:   input.ID,
			Done:     false,
			Position: 0,
		},
		{
			Title:    "Welcome to Godo! (Home example)",
			Category: "Home",
			Priority: 2,
			UserID:   input.ID,
			Done:     false,
			Position: 0,
		},
		{
			Title:    "Organize my workspace (Work example)",
			Category: "Work",
			Priority: 1,
			UserID:   input.ID,
			Done:     false,
			Position: 1,
		},
	}

	// Batch insert the defaults
	database.DB.Create(&defaultTasks)

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(config.JwtSecret)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// AuthMiddleware needs to be here to access JwtKey
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}
		tokenString := strings.Split(authHeader, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return config.JwtSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["user_id"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
	}
}
