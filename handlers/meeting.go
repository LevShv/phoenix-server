package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"phoenix-server/db"
	"phoenix-server/models"
	"time"

	"github.com/gin-gonic/gin"
)

func generateMeetingID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func CreateMeeting(c *gin.Context) {
	userType := c.GetString("user_type")
	if userType != "speaker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только выступающие могут создавать встречи"})
		return
	}

	userID := c.GetInt("user_id")

	var req models.CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	meetingID := generateMeetingID()
	now := time.Now().Unix()

	_, err := db.DB.Exec(`
        INSERT INTO meetings (id, title, description, location, start_date, speaker_id, created_at, status)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, meetingID, req.Title, req.Description, req.Location, req.StartDate, userID, now, "active")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания встречи"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Встреча создана",
		"meeting_id": meetingID,
	})
}

func GetAllMeetings(c *gin.Context) {
	rows, err := db.DB.Query(`
        SELECT m.id, m.title, m.description, m.location, m.start_date, 
               m.speaker_id, m.created_at, m.status, u.username
        FROM meetings m
        JOIN users u ON m.speaker_id = u.id
        ORDER BY m.created_at DESC
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения встреч"})
		return
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		var speakerName string
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Location, &m.StartDate,
			&m.SpeakerID, &m.CreatedAt, &m.Status, &speakerName)
		if err != nil {
			continue
		}
		meetings = append(meetings, m)
	}

	c.JSON(http.StatusOK, meetings)
}

func GetMyMeetings(c *gin.Context) {
	userType := c.GetString("user_type")
	if userType != "speaker" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступно только для выступающих"})
		return
	}

	userID := c.GetInt("user_id")

	rows, err := db.DB.Query(`
        SELECT id, title, description, location, start_date, created_at, status
        FROM meetings
        WHERE speaker_id = ?
        ORDER BY created_at DESC
    `, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения встреч"})
		return
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.Location, &m.StartDate, &m.CreatedAt, &m.Status)
		if err != nil {
			continue
		}
		m.SpeakerID = userID
		meetings = append(meetings, m)
	}

	c.JSON(http.StatusOK, meetings)
}

func GetMeetingByID(c *gin.Context) {
	meetingID := c.Param("id")

	var m models.Meeting
	err := db.DB.QueryRow(`
        SELECT id, title, description, location, start_date, speaker_id, created_at, status
        FROM meetings
        WHERE id = ?
    `, meetingID).Scan(&m.ID, &m.Title, &m.Description, &m.Location, &m.StartDate, &m.SpeakerID, &m.CreatedAt, &m.Status)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	c.JSON(http.StatusOK, m)
}
