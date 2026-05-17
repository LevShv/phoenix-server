package models

type Meeting struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	StartDate   string `json:"start_date"`
	SpeakerID   int    `json:"speaker_id"`
	CreatedAt   int64  `json:"created_at"`
	Status      string `json:"status"`
}

type CreateMeetingRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Location    string `json:"location" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	Status      string `json:"status"`
}
