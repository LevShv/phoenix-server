package models

type Poll struct {
	ID        int    `json:"id"`
	MeetingID string `json:"meeting_id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	CreatedAt int64  `json:"created_at"`
	QR        []byte `json:"-"`
}

type CreatePollRequest struct {
	Title string `json:"title" binding:"required"`
	URL   string `json:"url" binding:"required"`
}
