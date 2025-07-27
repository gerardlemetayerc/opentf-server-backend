package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

func GetProperty(c *gin.Context) {
	db := db.GetDB()
	propertyID := c.Param("property_id")
	var property models.Property
	if db.First(&property, propertyID).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}
	c.JSON(http.StatusOK, property)
}

func UpdateProperty(c *gin.Context) {
	db := db.GetDB()
	propertyID := c.Param("property_id")
	var property models.Property
	if db.First(&property, propertyID).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&property)
	c.JSON(http.StatusOK, property)
}

func DeleteProperty(c *gin.Context) {
	db := db.GetDB()
	propertyID := c.Param("property_id")
	var property models.Property
	if db.First(&property, propertyID).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}
	db.Delete(&property)
	c.JSON(http.StatusNoContent, nil)
}

func ListProperties(c *gin.Context) {
	moduleID := c.Param("module_id")
	db := db.GetDB()
	var properties []models.Property
	db.Where("module_id = ?", moduleID).Find(&properties)
	c.JSON(http.StatusOK, properties)
}

func CreateProperty(c *gin.Context) {
	moduleID := c.Param("module_id")
	db := db.GetDB()
	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	property.ModuleID = parseUint(moduleID)
	db.Create(&property)
	c.JSON(http.StatusCreated, property)
}

func parseUint(s string) uint {
	var u uint
	fmt.Sscanf(s, "%d", &u)
	return u
}
