package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

// --- OfferCategory CRUD ---
func GetOfferCategories(c *gin.Context) {
	db := db.GetDB()
	var cats []models.OfferCategory
	db.Find(&cats)
	c.JSON(http.StatusOK, cats)
}

func CreateOfferCategory(c *gin.Context) {
	db := db.GetDB()
	var input models.OfferCategory
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

func UpdateOfferCategory(c *gin.Context) {
	db := db.GetDB()
	var cat models.OfferCategory
	id := c.Param("id")
	if err := db.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	var input models.OfferCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.Name = input.Name
	if err := db.Save(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func DeleteOfferCategory(c *gin.Context) {
	db := db.GetDB()
	var cat models.OfferCategory
	id := c.Param("id")
	if err := db.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if err := db.Delete(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// --- Offer CRUD ---
func GetOffers(c *gin.Context) {
	db := db.GetDB()
	var offers []models.Offer
	db.Preload("Category").Find(&offers)
	c.JSON(http.StatusOK, offers)
}

func GetOffer(c *gin.Context) {
	db := db.GetDB()
	var offer models.Offer
	offerID := c.Param("offer_id")
	if err := db.Preload("Category").First(&offer, offerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
		return
	}
	c.JSON(http.StatusOK, offer)
}

func CreateOffer(c *gin.Context) {
	db := db.GetDB()
	var input models.Offer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Preload("Category").First(&input, input.ID)
	c.JSON(http.StatusOK, input)
}

func UpdateOffer(c *gin.Context) {
	db := db.GetDB()
	var offer models.Offer
	offerID := c.Param("offer_id")
	if err := db.First(&offer, offerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
		return
	}
	var input models.Offer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offer.Name = input.Name
	offer.Version = input.Version
	offer.Icon = input.Icon
	offer.GitURL = input.GitURL
	offer.Active = input.Active
	offer.CategoryID = input.CategoryID
	if err := db.Save(&offer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Preload("Category").First(&offer, offer.ID)
	c.JSON(http.StatusOK, offer)
}

func DeleteOffer(c *gin.Context) {
	db := db.GetDB()
	var offer models.Offer
	offerID := c.Param("offer_id")
	if err := db.First(&offer, offerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
		return
	}
	if err := db.Delete(&offer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
