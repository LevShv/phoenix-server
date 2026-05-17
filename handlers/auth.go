package handlers

import (
	_ "log"
	"net/http"
	"phoenix-server/db"
	"phoenix-server/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("phoenix")

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	_, err = db.DB.Exec(`
        INSERT INTO users (username, password_hash, name, surname, patronymic, type)
        VALUES (?, ?, ?, ?, ?, ?)
    `, req.Username, string(hash), req.Name, req.Surname, req.Patronymic, req.Type)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь уже существует"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь создан"})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	var user models.User
	var hash string
	err := db.DB.QueryRow(`
        SELECT id, username, password_hash, name, surname, patronymic, type 
        FROM users 
        WHERE username = ?
    `, req.Username).Scan(
		&user.ID, &user.Username, &hash,
		&user.Name, &user.Surname, &user.Patronymic, &user.Type,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные имя пользователя или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные имя пользователя или пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"type":     user.Type,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: tokenString})
}

func GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	var user models.User
	err := db.DB.QueryRow(`
        SELECT id, username, name, surname, patronymic, type
        FROM users 
        WHERE id = ?
    `, userID).Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Patronymic, &user.Type)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения профиля"})
		return
	}

	c.JSON(http.StatusOK, user)
}
