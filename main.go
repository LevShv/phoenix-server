package main

import (
	_ "log"
	"phoenix-server/db"
	//	"github.com/gin-gonic/gin"
	//
	// "phoenix-server/handlers"
	// "phoenix-server/middleware"
)

func main() {
	db.InitDB()
	defer db.Close()

	//r := gin.Default();

	//r.POST(/"register")

}
