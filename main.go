package main

import (
	_ "log"
	"phoenix-server/db"

	"github.com/gin-gonic/gin"

	"phoenix-server/handlers"
	"phoenix-server/middleware"
)

func main() {
	db.InitDB()
	defer db.Close()

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/profile", handlers.GetProfile)

		api.POST("/meetings", handlers.CreateMeeting)
		api.GET("/meetings", handlers.GetAllMeetings)
		api.GET("/meetings/my", handlers.GetMyMeetings)
		api.GET("/meetings/:id", handlers.GetMeetingByID)

		api.POST("/meetings/:id/polls", handlers.AddPoll)
		api.GET("/meetings/:id/polls", handlers.GetMeetingPolls)
		api.DELETE("/meetings/:id/polls/:poll_id", handlers.DeletePoll)

	}

	r.Run(":8080")
}
