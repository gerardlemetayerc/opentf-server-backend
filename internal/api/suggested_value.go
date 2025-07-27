package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

// CRUD Domain
func GetDomains(c *gin.Context) {
	db := db.GetDB()
	var domains []models.Domain
	db.Find(&domains)
	c.JSON(http.StatusOK, domains)
}

func GetDomain(c *gin.Context) {
	db := db.GetDB()
	var domain models.Domain
	id := c.Param("domain_id")
	if err := db.First(&domain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}
	c.JSON(http.StatusOK, domain)
}

func CreateDomain(c *gin.Context) {
	db := db.GetDB()
	var input models.Domain
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

func UpdateDomain(c *gin.Context) {
	db := db.GetDB()
	var domain models.Domain
	id := c.Param("domain_id")
	if err := db.First(&domain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}
	var input models.Domain
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	domain.Name = input.Name
	domain.Label = input.Label
	if err := db.Save(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain)
}

func DeleteDomain(c *gin.Context) {
	db := db.GetDB()
	var domain models.Domain
	id := c.Param("domain_id")
	if err := db.First(&domain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}
	if err := db.Delete(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// CRUD SuggestedValue
func GetSuggestedValues(c *gin.Context) {
	db := db.GetDB()
	var values []models.SuggestedValue

	query := db.Preload("Domain")

	// Filtrage par domain_id
	if domainID := c.Query("domain_id"); domainID != "" {
		query = query.Where("domain_id = ?", domainID)
	}
	// Filtrage par parent_domain_id
	if parentDomainID := c.Query("parent_domain_id"); parentDomainID != "" {
		query = query.Where("parent_domain_id = ?", parentDomainID)
	}
	// Filtrage par parent_value
	if parentValue := c.Query("parent_value"); parentValue != "" {
		query = query.Where("parent_value = ?", parentValue)
	}

	query.Find(&values)
	c.JSON(http.StatusOK, values)
}

func GetSuggestedValue(c *gin.Context) {
	db := db.GetDB()
	var value models.SuggestedValue
	id := c.Param("id")
	if err := db.Preload("Domain").First(&value, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SuggestedValue not found"})
		return
	}
	c.JSON(http.StatusOK, value)
}

func CreateSuggestedValue(c *gin.Context) {
	db := db.GetDB()
	var input models.SuggestedValue
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Preload("Domain").First(&input, input.ID)
	c.JSON(http.StatusOK, input)
}

func UpdateSuggestedValue(c *gin.Context) {
	db := db.GetDB()
	var value models.SuggestedValue
	id := c.Param("id")
	if err := db.First(&value, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SuggestedValue not found"})
		return
	}
	var input models.SuggestedValue
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value.DomainID = input.DomainID
	value.DisplayValue = input.DisplayValue
	value.RealValue = input.RealValue
	if err := db.Save(&value).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Preload("Domain").First(&value, value.ID)
	c.JSON(http.StatusOK, value)
}

func DeleteSuggestedValue(c *gin.Context) {
	db := db.GetDB()
	var value models.SuggestedValue
	id := c.Param("id")
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
