package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-gin-jwt-demo/config"
	"go-gin-jwt-demo/middleware"
	"go-gin-jwt-demo/model"

	"log/slog"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required"`
}

// Register creates a user account.
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Enter a username of 1–50 characters and a password.",
		})
		return
	}

	if len(req.Password) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Password must not exceed 72 bytes.",
		})
		return
	}

	hashPwd, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Unable to process registration.",
		})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashPwd),
	}

	err = config.DB.Create(&user).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "This username is already taken. Choose another or log in.",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Unable to save your account. Please try again.",
		})
		return
	}

	slog.Info("registration_success",
		"username", user.Username,
		"client_ip", c.ClientIP(),
	)

	c.JSON(http.StatusCreated, gin.H{
		"msg": "Registration successful.",
	})
}

// Login checks credentials and returns a JWT.
func Login(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Enter a username of 1–50 characters and a password.",
		})
		return
	}

	var user model.User
	err := config.DB.Where("username = ?", req.Username).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Username or password is incorrect.",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Login is unavailable. Please try again.",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {

		slog.Error("login_failed",
			"username", req.Username,
			"client_ip", c.ClientIP(),
			"reason", "password_verification_failed",
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Username or password is incorrect.",
		})
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Unable to complete login. Please try again.",
		})
		return
	}

	slog.Info("login_success",
		"username", user.Username,
		"client_ip", c.ClientIP(),
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Info returns the logged-in user's profile.
func Info(c *gin.Context) {
	uid, exists := c.Get("userId")
	userID, ok := uid.(uint)

	if !exists || !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Please log in again.",
		})
		return
	}

	var user model.User
	result := config.DB.First(&user, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Please log in again.",
		})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Unable to load your profile.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":    user.ID,
		"username":  user.Username,
		"createdAt": user.CreatedAt,
		"msg":       "Profile loaded successfully.",
	})
}
