package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

// GET /backendapi/statefiles/:id
func GetStateFile(c *gin.Context) {
	db := db.GetDB()
	var state models.StateFile
	id := c.Param("id")
	if err := db.First(&state, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "StateFile not found"})
		return
	}
	c.JSON(http.StatusOK, state)
}

// POST /backendapi/statefiles
func CreateOrUpdateStateFile(c *gin.Context) {
	db := db.GetDB()
	var input struct {
		InstanceID uint   `json:"instance_id"`
		Data       []byte `json:"data"`
		Version    int    `json:"version"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var state models.StateFile
	if err := db.Where("instance_id = ?", input.InstanceID).Last(&state).Error; err == nil {
		// Update
		state.Data = input.Data
		state.Version = input.Version
		if err := db.Save(&state).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, state)
		return
	}
	// Create
	state = models.StateFile{
		InstanceID: input.InstanceID,
		Data:       input.Data,
		Version:    input.Version,
	}
	if err := db.Create(&state).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, state)
}
