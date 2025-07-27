package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
	"strconv"
)

// GET /api/offers/:offer_id/properties
func GetOfferProperties(c *gin.Context) {
	db := db.GetDB()
	offerID, _ := strconv.Atoi(c.Param("offer_id"))
	var props []models.OfferProperty
	db.Where("offer_id = ?", offerID).Find(&props)
	c.JSON(http.StatusOK, props)
}

// GET /api/offers/:offer_id/properties/:id
func GetOfferProperty(c *gin.Context) {
	db := db.GetDB()
	var prop models.OfferProperty
	id := c.Param("id")
	if err := db.First(&prop, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}
	c.JSON(http.StatusOK, prop)
}

// POST /api/offers/:offer_id/properties
func CreateOfferProperty(c *gin.Context) {
	db := db.GetDB()
	offerID, _ := strconv.Atoi(c.Param("offer_id"))
	var input models.OfferProperty
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.OfferID = uint(offerID)
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

// PUT /api/offers/:offer_id/properties/:id
func UpdateOfferProperty(c *gin.Context) {
	db := db.GetDB()
	var prop models.OfferProperty
	id := c.Param("id")
	if err := db.First(&prop, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}
	var input models.OfferProperty
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	prop.Name = input.Name
	prop.Type = input.Type
	prop.Label = input.Label
	prop.Description = input.Description
	prop.Required = input.Required
	prop.DefaultValue = input.DefaultValue
	prop.MinValue = input.MinValue
	prop.MaxValue = input.MaxValue
	prop.MetadataSource = input.MetadataSource
	prop.DependsOn = input.DependsOn
	prop.CustomJS = input.CustomJS
	if err := db.Save(&prop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prop)
}

// DELETE /api/offers/:offer_id/properties/:id
func DeleteOfferProperty(c *gin.Context) {
	db := db.GetDB()
	var prop models.OfferProperty
	id := c.Param("id")
	if err := db.First(&prop, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}
	if err := db.Delete(&prop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
