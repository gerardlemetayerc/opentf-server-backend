package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

func ListInstances(c *gin.Context) {
	db := db.GetDB()
	var instances []models.Instance
	db.Find(&instances)
	c.JSON(http.StatusOK, instances)
}
