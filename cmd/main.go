package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"opentf-server/internal/api"
	"os"
	"strings"
)

func main() {
	r := gin.Default()
	origins := os.Getenv("FRONTEND_ORIGINS")
	originList := strings.Split(origins, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     originList,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	api.RegisterRoutes(r)
	r.Run(":8080")
}
