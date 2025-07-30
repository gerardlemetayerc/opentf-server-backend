package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
	"time"
)

// POST /backendapi/statelocks
func AcquireStateLock(c *gin.Context) {
	db := db.GetDB()
	var input struct {
		InstanceID uint   `json:"instance_id"`
		Info       string `json:"info"`
		TTL        int    `json:"ttl"` // en secondes
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Vérifie s'il existe déjà un lock actif
	var lock models.StateLock
	if err := db.Where("instance_id = ? AND expires_at > ?", input.InstanceID, time.Now()).First(&lock).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Lock already acquired", "lock": lock})
		return
	}
	// Crée le lock
	lock = models.StateLock{
		InstanceID: input.InstanceID,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(time.Duration(input.TTL) * time.Second),
		Info:       input.Info,
	}
	if err := db.Create(&lock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lock)
}

// DELETE /backendapi/statelocks/:id
func ReleaseStateLock(c *gin.Context) {
	db := db.GetDB()
	id := c.Param("id")
	if err := db.Delete(&models.StateLock{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"released": true})
}

// GET /backendapi/statelocks/:instance_id
func GetStateLock(c *gin.Context) {
	db := db.GetDB()
	instanceID := c.Param("instance_id")
	var lock models.StateLock
	if err := db.Where("instance_id = ? AND expires_at > ?", instanceID, time.Now()).First(&lock).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active lock"})
		return
	}
	c.JSON(http.StatusOK, lock)
}
