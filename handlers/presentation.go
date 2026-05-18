package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"phoenix-server/db"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadPresentation godoc
// @Summary      Загрузить презентацию
// @Security     BearerAuth
// @Param        id path string true "ID встречи"
// @Param        presentation formData file true "PDF файл"
// @Router       /api/meetings/{id}/presentation [post]
func UploadPresentation(c *gin.Context) {
	meetingID := c.Param("id")
	userID := c.GetInt("user_id")

	var speakerID int
	var oldURL *string
	err := db.DB.QueryRow("SELECT speaker_id, presentation_url FROM meetings WHERE id = ?", meetingID).
		Scan(&speakerID, &oldURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if speakerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только создатель может загружать презентации"})
		return
	}

	file, err := c.FormFile("presentation")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}

	ext := filepath.Ext(file.Filename)
	allowedExt := map[string]bool{
		".pdf":  true,
		".pptx": true,
		".ppt":  true,
	}
	if !allowedExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Разрешены только PDF, PPT, PPTX"})
		return
	}

	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл слишком большой (макс 50MB)"})
		return
	}

	if oldURL != nil && *oldURL != "" {
		oldPath := "./" + *oldURL
		os.Remove(oldPath)
	}

	timestamp := time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%s_%d%s", meetingID, timestamp, ext)
	filePath := filepath.Join("uploads/meetings/presentations/", uniqueFilename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}

	presentationURL := "uploads/meetings/presentations/" + uniqueFilename
	_, err = db.DB.Exec("UPDATE meetings SET presentation_url = ? WHERE id = ?", presentationURL, meetingID)
	if err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления БД"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Презентация загружена",
		"url":     presentationURL,
		"size":    file.Size,
	})
}

func GetPresentation(c *gin.Context) {
	meetingID := c.Param("id")

	var presentationURL *string
	err := db.DB.QueryRow("SELECT presentation_url FROM meetings WHERE id = ?", meetingID).Scan(&presentationURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if presentationURL == nil || *presentationURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Презентация не загружена"})
		return
	}

	filePath := "./" + *presentationURL
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		db.DB.Exec("UPDATE meetings SET presentation_url = NULL WHERE id = ?", meetingID)
		c.JSON(http.StatusOK, gin.H{"has_presentation": false})
		return
	}

	c.File(filePath)
}

func GetPresentationInfo(c *gin.Context) {
	meetingID := c.Param("id")

	var presentationURL *string
	err := db.DB.QueryRow("SELECT presentation_url FROM meetings WHERE id = ?", meetingID).Scan(&presentationURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if presentationURL == nil || *presentationURL == "" {
		c.JSON(http.StatusOK, gin.H{"has_presentation": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_presentation": true,
		"url":              *presentationURL,
	})
}

func DeletePresentation(c *gin.Context) {
	meetingID := c.Param("id")
	userID := c.GetInt("user_id")

	var speakerID int
	var presentationURL *string
	err := db.DB.QueryRow("SELECT speaker_id, presentation_url FROM meetings WHERE id = ?", meetingID).
		Scan(&speakerID, &presentationURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if speakerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только создатель может удалять презентации"})
		return
	}

	if presentationURL == nil || *presentationURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Презентация не найдена"})
		return
	}

	filePath := "./" + *presentationURL
	os.Remove(filePath)

	_, err = db.DB.Exec("UPDATE meetings SET presentation_url = NULL WHERE id = ?", meetingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Презентация удалена"})
}
