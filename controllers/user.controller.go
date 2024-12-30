package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mmycin/ndc14/config"
	"github.com/mmycin/ndc14/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var Body struct {
		Username string `json:"username"`
		Password string `json:"-"`
		Email    string `json:"email" gorm:"unique"`
		Roll     string `json:"roll" gorm:"unique"`
		Batch    int    `json:"batch"`
		FBLink   string `json:"fbLink"`
		IsAdmin  bool   `json:"isAdmin" gorm:"default:false"`
	}

	if err := c.ShouldBindJSON(&Body); err != nil {
		c.JSON(400, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user
	user := models.User{
		Username: Body.Username,
		Email:    Body.Email,
		Password: string(hash),
		Roll:     Body.Roll,
		Batch:    Body.Batch,
		FBLink:   Body.FBLink,
		IsAdmin:  false,
	}
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Error creating User",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"data":    user,
	})
}

func Login(c *gin.Context) {
	var Body struct {
		Email    string `json:"email"`
		Password string `json:"-"`
	}

	if err := c.ShouldBindJSON(&Body); err != nil {
		c.JSON(400, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	var user models.User
	config.DB.First(&user, "email = ?", Body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid Email",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Body.Password))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Invalid Password, try again",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status": "An error occured while creating the Token",
			"error":  err,
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": "You are logged in",
		"data":    user,
	})
}

func Logout(c *gin.Context) {
	// Clear the Authorization cookie by setting its max age to -1
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}
