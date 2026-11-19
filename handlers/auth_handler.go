package handlers

import (
	"net/http"
	"pr-5-go/database"
	"pr-5-go/dto"
	"pr-5-go/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("super_secret_key_for_incident_tracker")


func Register(c *gin.Context) {
	var input dto.RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(err)})
		return
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка хешування пароля"})
		return
	}

	user := models.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Користувач із таким email вже існує"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Користувача успішно зареєстровано",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var input dto.LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(err)})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невірний email або пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невірний email або пароль"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Помилка генерації токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}