package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"opentf-server/internal/api"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	api.RegisterRoutes(r)
	r.Run(":8080")
}
