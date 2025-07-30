package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

// GET /api/providers
func GetProviders(c *gin.Context) {
	db := db.GetDB()
	var providers []models.Provider
	db.Find(&providers)
	c.JSON(http.StatusOK, providers)
}

// POST /api/providers
func CreateProvider(c *gin.Context) {
	db := db.GetDB()
	var input models.Provider
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}
