package main

import (
	_ "log"
	"phoenix-server/db"

	"github.com/gin-gonic/gin"

	"phoenix-server/handlers"
	"phoenix-server/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "phoenix-server/docs"
)

func main() {
	db.InitDB()
	defer db.Close()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

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

		api.POST("/meetings/:id/presentation", handlers.UploadPresentation)
		api.GET("/meetings/:id/presentation", handlers.GetPresentation)
		api.GET("/meetings/:id/presentation/info", handlers.GetPresentationInfo)
		api.DELETE("/meetings/:id/presentation", handlers.DeletePresentation)

		api.POST("/meetings/:id/photo", handlers.UploadMeetingPhoto)
		api.GET("/meetings/:id/photo", handlers.GetMeetingPhoto)
		api.GET("/meetings/:id/photo/info", handlers.GetMeetingPhotoInfo)
		api.DELETE("/meetings/:id/photo", handlers.DeleteMeetingPhoto)

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run("0.0.0.0:8080")
}
