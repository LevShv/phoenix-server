package handlers

import (
	"net/http"
	"phoenix-server/db"
	"phoenix-server/models"
	"time"

	"github.com/gin-gonic/gin"
)

func AddPoll(c *gin.Context) {
	meetingID := c.Param("id")
	userID := c.GetInt("user_id")

	var speakerID int
	err := db.DB.QueryRow("SELECT speaker_id FROM meetings WHERE id = ?", meetingID).Scan(&speakerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if speakerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Только создатель встречи может добавлять опросы"})
		return
	}

	var req models.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	_, err = db.DB.Exec(`
        INSERT INTO polls (meeting_id, title, url, created_at)
        VALUES (?, ?, ?, ?)
    `, meetingID, req.Title, req.URL, time.Now().Unix())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания опроса"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Опрос добавлен"})
}

func GetMeetingPolls(c *gin.Context) {
	meetingID := c.Param("id")

	rows, err := db.DB.Query(`
        SELECT id, meeting_id, title, url, created_at
        FROM polls
        WHERE meeting_id = ?
        ORDER BY created_at DESC
    `, meetingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения опросов"})
		return
	}
	defer rows.Close()

	var polls []models.Poll
	for rows.Next() {
		var p models.Poll
		err := rows.Scan(&p.ID, &p.MeetingID, &p.Title, &p.URL, &p.CreatedAt)
		if err != nil {
			continue
		}
		polls = append(polls, p)
	}

	c.JSON(http.StatusOK, polls)
}

func DeletePoll(c *gin.Context) {
	meetingID := c.Param("id")
	pollID := c.Param("poll_id")
	userID := c.GetInt("user_id")

	var speakerID int
	err := db.DB.QueryRow("SELECT speaker_id FROM meetings WHERE id = ?", meetingID).Scan(&speakerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Встреча не найдена"})
		return
	}

	if speakerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM polls WHERE id = ? AND meeting_id = ?", pollID, meetingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Опрос удалён"})
}
