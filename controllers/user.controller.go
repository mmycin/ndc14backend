package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mmycin/ndc14/config"
	"github.com/mmycin/ndc14/libs"
	"github.com/mmycin/ndc14/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var Body struct {
		FullName string `json:"fullName" gorm:"not null"`
		Username string `json:"username" gorm:"not null"`
		Password string `json:"password" gorm:"not null"`
		Email    string `json:"email" gorm:"unique;not null"`
		Roll     string `json:"roll" gorm:"unique;not null"`
		Batch    int    `json:"batch" gorm:"not null"`
		FBLink   string `json:"fbLink" gorm:"not null"`
		IsAdmin  bool   `json:"isAdmin" gorm:"default:false;not null"`
	}

	if err := c.ShouldBindJSON(&Body); err != nil {
		c.JSON(400, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	// Trim any whitespace from password
	// Body.Password = strings.TrimSpace(Body.Password)

	if len(Body.Password) < 1 {
		c.JSON(400, gin.H{
			"error": "Password cannot be empty",
		})
		return
	}

	// Hash password with specific cost
	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), bcrypt.DefaultCost)
	fmt.Println(hash)
	if err != nil {
		fmt.Printf("Password hashing error: %v\n", err)
		c.JSON(400, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user with trimmed and properly hashed password
	user := models.User{
		FullName: Body.FullName,
		Username: Body.Username,
		Email:    Body.Email,
		// Password: string(hash),
		Password: string(hash),
		Roll:     Body.Roll,
		Batch:    Body.Batch,
		FBLink:   Body.FBLink,
		IsAdmin:  false,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		fmt.Printf("Database error during user creation: %v\n", result.Error)
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Error creating User",
		})
		return
	}

	// Log success but don't expose hash in response
	fmt.Printf("Successfully created user with roll: %s\n", user.Roll)

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"data":    gin.H{
			"id":   user.ID,
			"roll": user.Roll,
			"email": user.Email,
			"fbLink": user.FBLink,
			"isAdmin": user.IsAdmin,
			"fullName": user.FullName,
			"username": user.Username,
			"batch": user.Batch,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
			"deletedAt": user.DeletedAt,
		},
	})
}

func Login(c *gin.Context) {
	var Body struct {
		Username string `json:"username" gorm:"not null"`
		Roll     string `json:"roll" gorm:"unique;not null"`
		Password string `json:"password" gorm:"not null"`
	}

	if err := c.ShouldBindJSON(&Body); err != nil {
		c.JSON(400, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	// Trim any whitespace from password
	Body.Password = strings.TrimSpace(Body.Password)

	var user models.User
	if err := config.DB.Where("roll = ?", Body.Roll).First(&user).Error; err != nil {
		fmt.Printf("User lookup error for roll %s: %v\n", Body.Roll, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Add logging to check password lengths and content
	fmt.Printf("Login attempt for roll: %s\n", Body.Roll)
	fmt.Printf("Stored hash length: %d\n", len(user.Password))
	fmt.Printf("Provided password length: %d\n", len(Body.Password))

	// Compare hashed password with the entered password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Body.Password))
	// err := strings.Compare(Body.Password, user.Password)
	if err != nil {
		fmt.Printf("Password comparison error for roll %s: %v\n", Body.Roll, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Error creating token",
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

func UpdateUser(c *gin.Context) {
	user, _ := c.Get("user")
	var existingUser models.User
	config.DB.First(&existingUser, "id = ?", user.(models.User).ID)

	if existingUser.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the user is admin or is trying to update their own account
	if existingUser.IsAdmin || existingUser.ID == user.(models.User).ID {
		var Body struct {
			FullName string `json:"fullName"`
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email" gorm:"unique"`
			Roll     string `json:"roll" gorm:"unique"`
			Batch    int    `json:"batch"`
			FBLink   string `json:"fbLink"`
		}

		if err := c.ShouldBindJSON(&Body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if Body.FullName != "" {
			existingUser.FullName = Body.FullName
		}
		if Body.Username != "" {
			existingUser.Username = Body.Username
		}
		if Body.Email != "" {
			if !libs.IsValidEmail(Body.Email) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid email format",
				})
				return
			}
			existingUser.Email = Body.Email
		}
		if Body.FBLink != "" {
			existingUser.FBLink = Body.FBLink
		}
		if Body.Roll != "" {
			if !libs.IsValidRoll(Body.Roll) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid roll format. Must follow pattern: 1(2|3)x14xxx where x is any digit and last 3 digits between 001-150",
				})
				return
			}
			existingUser.Roll = Body.Roll
		}
		if Body.Batch != 0 {
			existingUser.Batch = Body.Batch
		}
		if Body.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
				return
			}
			existingUser.Password = string(hash)
		}

		config.DB.Save(&existingUser)
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "data": existingUser})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can only update your own account"})
	}
}

func DeleteUser(c *gin.Context) {
	user, _ := c.Get("user")
	var existingUser models.User
	var Body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Roll     string `json:"roll"`
	}

	if err := c.ShouldBindJSON(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Admin flow - can delete any user including themselves
	if user.(models.User).IsAdmin {
		if Body.Username != "" {
			if result := config.DB.Delete(&existingUser, "username = ?", Body.Username); result.Error != nil || result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
		} else if Body.Email != "" {
			if result := config.DB.Delete(&existingUser, "email = ?", Body.Email); result.Error != nil || result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
		} else if Body.Roll != "" {
			if result := config.DB.Delete(&existingUser, "roll = ?", Body.Roll); result.Error != nil || result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide username, email, or roll to delete user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
		return
	}

	// Regular user flow - can only delete their own account
	config.DB.Delete(&existingUser, "id = ?", user.(models.User).ID)
	c.JSON(http.StatusOK, gin.H{"message": "Your account has been deleted successfully"})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func GetUserByUsername(c *gin.Context) {
	var user models.User
	username := c.Param("username")

	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func GetUserByRoll(c *gin.Context) {
	var user models.User
	roll := c.Param("roll")

	if err := config.DB.Where("roll = ?", roll).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
