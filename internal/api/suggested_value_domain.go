package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
	"strconv"
)

// GET /api/domains/:domain_id/suggested_values
func GetSuggestedValuesByDomain(c *gin.Context) {
	db := db.GetDB()
	domainID, err := strconv.Atoi(c.Param("domain_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain_id"})
		return
	}
	var values []models.SuggestedValue
	db.Where("domain_id = ?", domainID).Preload("Domain").Find(&values)
	c.JSON(http.StatusOK, values)
}

// POST /api/domains/:domain_id/suggested_values
func CreateSuggestedValueByDomain(c *gin.Context) {
	db := db.GetDB()
	domainID, err := strconv.Atoi(c.Param("domain_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain_id"})
		return
	}
	var input models.SuggestedValue
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.DomainID = uint(domainID)
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Preload("Domain").First(&input, input.ID)
	c.JSON(http.StatusOK, input)
}

// PUT /api/domains/:domain_id/suggested_values/:id
func UpdateSuggestedValueByDomain(c *gin.Context) {
	db := db.GetDB()
	id := c.Param("id")
	var value models.SuggestedValue
	if err := db.First(&value, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SuggestedValue not found"})
		return
	}
	var input models.SuggestedValue
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value.DisplayValue = input.DisplayValue
	value.RealValue = input.RealValue
	value.ParentDomainID = input.ParentDomainID
	value.ParentValue = input.ParentValue
	if err := db.Save(&value).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Preload("Domain").First(&value, value.ID)
	c.JSON(http.StatusOK, value)
}

// DELETE /api/domains/:domain_id/suggested_values/:id
func DeleteSuggestedValueByDomain(c *gin.Context) {
	db := db.GetDB()
	id := c.Param("id")
	var value models.SuggestedValue
	if err := db.First(&value, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SuggestedValue not found"})
		return
	}
	if err := db.Delete(&value).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
