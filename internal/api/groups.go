package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

// GET /api/groups
func GetGroups(c *gin.Context) {
	db := db.GetDB()
	var groups []models.Group
	db.Find(&groups)
	c.JSON(http.StatusOK, groups)
}

// POST /api/groups
func CreateGroup(c *gin.Context) {
	db := db.GetDB()
	var input models.Group
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

// PUT /api/groups/:id
func UpdateGroup(c *gin.Context) {
	db := db.GetDB()
	var group models.Group
	id := c.Param("id")
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	var input models.Group
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	group.Name = input.Name
	group.Roles = input.Roles
	if err := db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

// DELETE /api/groups/:id
func DeleteGroup(c *gin.Context) {
	db := db.GetDB()
	var group models.Group
	id := c.Param("id")
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	if err := db.Delete(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
