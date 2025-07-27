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

func CreateInstance(c *gin.Context) {
	db := db.GetDB()
	var instance models.Instance
	if err := c.ShouldBindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&instance)
	c.JSON(http.StatusCreated, instance)
}
