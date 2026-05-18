package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"phoenix-server/db"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	meetingphotopath = "uploads/meetings/photos/"
)

func UploadMeetingPhoto(c *gin.Context) {
	meetingID := c.Param("id")
	userID := c.GetInt("user_id")

	var speakerID int
	var oldURL *string
	err := db.DB.QueryRow("SELECT speaker_id, photo_url FROM meetings WHERE id = ?", meetingID).
		Scan(&speakerID, &oldURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if speakerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только создатель может загружать фото"})
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}

	ext := filepath.Ext(file.Filename)
	allowedExt := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jprg": true,
	}
	if !allowedExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимый формат файла"})
		return
	}

	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл слишком большой (макс 50MB)"})
		return
	}

	if oldURL != nil && *oldURL != "" {
		oldPath := "./" + *oldURL
		log.Printf("Path: %v", *oldURL)
		os.Remove(oldPath)
	}

	timestamp := time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%s_%d%s", meetingID, timestamp, ext)
	filePath := filepath.Join(meetingphotopath, uniqueFilename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}

	photoURL := meetingphotopath + uniqueFilename
	_, err = db.DB.Exec("UPDATE meetings SET photo_url = ? WHERE id = ?", photoURL, meetingID)
	if err != nil {
		log.Printf(err.Error())
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления БД"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Фото встречи загружено",
		"url":     photoURL,
		"size":    file.Size,
	})
}

func GetMeetingPhoto(c *gin.Context) {
	meetingID := c.Param("id")

	var photoURL *string
	err := db.DB.QueryRow("SELECT photo_url FROM meetings WHERE id = ?", meetingID).Scan(&photoURL)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	log.Printf("Path: %v", *photoURL)
	if photoURL == nil || *photoURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Картинка не загружена"})
		return
	}

	log.Printf("Path: %v", *photoURL)
	filePath := "./" + *photoURL
	c.File(filePath)
}

func DeleteMeetingPhoto(c *gin.Context) {
	meetingID := c.Param("id")
	userID := c.GetInt("user_id")

	var speakerID int
	var photoURL *string
	err := db.DB.QueryRow("SELECT speaker_id, photo_url FROM meetings WHERE id = ?", meetingID).
		Scan(&speakerID, &photoURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if speakerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только создатель может удалять фото"})
		return
	}

	if photoURL == nil || *photoURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Фото не найден"})
		return
	}

	filePath := "./" + *photoURL
	os.Remove(filePath)

	_, err = db.DB.Exec("UPDATE meetings SET photo_url = NULL WHERE id = ?", meetingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Фото удалено"})
}

func GetMeetingPhotoInfo(c *gin.Context) {
	meetingID := c.Param("id")

	var photoURL *string
	err := db.DB.QueryRow("SELECT photo_url FROM meetings WHERE id = ?", meetingID).Scan(&photoURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if photoURL == nil || *photoURL == "" {
		c.JSON(http.StatusOK, gin.H{"has_photo": false})
		return
	}

	filePath := "./" + *photoURL
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		db.DB.Exec("UPDATE meetings SET photo_url = NULL WHERE id = ?", meetingID)
		c.JSON(http.StatusOK, gin.H{"has_photo": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_photo": true,
		"url":       *photoURL,
	})
}
