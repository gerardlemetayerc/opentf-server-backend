package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

// Liste des instances
func GetInstances(c *gin.Context) {
	db := db.GetDB()
	var instances []models.Instance
	db.Preload("Properties").Find(&instances)

	// Récupère les infos avancées pour chaque instance
	var result []map[string]interface{}
	for _, inst := range instances {
		// Récupère le demandeur
		var requester models.User
		db.First(&requester, inst.RequesterID)
		displayName := requester.FirstName + " " + requester.LastName
		// Récupère l'offre
		var offer models.Offer
		db.First(&offer, inst.OfferID)
		// Utilise le champ Name si valorisé, sinon génère un identifiant arbitraire
		instanceName := inst.Name
		if instanceName == "" {
			instanceName = "instance-" + fmt.Sprint(inst.ID)
		}
		result = append(result, map[string]interface{}{
			"id":                    inst.ID,
			"offer_id":              inst.OfferID,
			"offer_name":            offer.Name,
			"created_at":            inst.CreatedAt,
			"updated_at":            inst.UpdatedAt,
			"status":                inst.Status,
			"requester_id":          inst.RequesterID,
			"requester_displayname": displayName,
			"validator_id":          inst.ValidatorID,
			"properties":            inst.Properties,
			"instance_name":         instanceName,
		})
	}
	c.JSON(http.StatusOK, result)
}

// Détail d'une instance
func GetInstance(c *gin.Context) {
	db := db.GetDB()
	var instance models.Instance
	id := c.Param("id")
	if err := db.Preload("Properties").First(&instance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}
	// Enrichit chaque propriété avec son nom
	var enrichedProps []map[string]interface{}
	for _, prop := range instance.Properties {
		var offerProp models.Property
		db.First(&offerProp, prop.OfferPropertyID)
		enrichedProps = append(enrichedProps, map[string]interface{}{
			"id":                prop.ID,
			"offer_property_id": prop.OfferPropertyID,
			"name":              offerProp.Name,
			"value":             prop.Value,
		})
	}
	// Ajoute les autres infos de l'instance
	var requester models.User
	db.First(&requester, instance.RequesterID)
	displayName := requester.FirstName + " " + requester.LastName
	var offer models.Offer
	db.First(&offer, instance.OfferID)
	instanceName := "instance-" + fmt.Sprint(instance.ID)
	// Ajoute les infos générales de l'offre
	offerInfo := map[string]interface{}{
		"id":                  offer.ID,
		"name":                offer.Name,
		"version":             offer.Version,
		"icon":                offer.Icon,
		"git_url":             offer.GitURL,
		"active":              offer.Active,
		"category_id":         offer.CategoryID,
		"auto_validated":      offer.AutoValidated,
		"validation_group_id": offer.ValidationGroupID,
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":                    instance.ID,
		"offer_id":              instance.OfferID,
		"offer_name":            offer.Name,
		"offer_info":            offerInfo,
		"created_at":            instance.CreatedAt,
		"updated_at":            instance.UpdatedAt,
		"status":                instance.Status,
		"requester_id":          instance.RequesterID,
		"requester_displayname": displayName,
		"validator_id":          instance.ValidatorID,
		"properties":            enrichedProps,
		"instance_name":         instanceName,
	})
}

// Création d'une instance liée à une offre
func CreateInstance(c *gin.Context) {
	db := db.GetDB()
	var input struct {
		OfferID     uint `json:"offer_id"`
		ValidatorID uint `json:"validator_id"`
		Properties  []struct {
			OfferPropertyID uint   `json:"offer_property_id"`
			Value           string `json:"value"`
		} `json:"properties"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Récupère l'utilisateur authentifié (exemple avec JWT)
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}
	// Récupère l'offre pour NamePropertyID
	var offer models.Offer
	if err := db.First(&offer, input.OfferID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offer not found"})
		return
	}
	var nameValue string
	for _, prop := range input.Properties {
		if offer.NamePropertyID != nil && prop.OfferPropertyID == *offer.NamePropertyID {
			nameValue = prop.Value
		}
	}
	instance := models.Instance{
		OfferID:     input.OfferID,
		Status:      "draft",
		RequesterID: userID.(uint),
		ValidatorID: input.ValidatorID,
		Name:        nameValue,
	}
	if err := db.Create(&instance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, prop := range input.Properties {
		instProp := models.InstanceProperty{
			InstanceID:      instance.ID,
			OfferPropertyID: prop.OfferPropertyID,
			Value:           prop.Value,
		}
		db.Create(&instProp)
	}
	db.Preload("Properties").First(&instance, instance.ID)
	c.JSON(http.StatusOK, instance)
}

// Mise à jour d'une instance
func UpdateInstance(c *gin.Context) {
	db := db.GetDB()
	var instance models.Instance
	id := c.Param("id")
	if err := db.First(&instance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}
	var input struct {
		Status      string `json:"status"`
		ValidatorID uint   `json:"validator_id"`
		Properties  []struct {
			OfferPropertyID uint   `json:"offer_property_id"`
			Value           string `json:"value"`
		} `json:"properties"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Status != "" {
		instance.Status = input.Status
	}
	if input.ValidatorID != 0 {
		instance.ValidatorID = input.ValidatorID
	}
	if err := db.Save(&instance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Mise à jour des propriétés et du nom si configuré
	var offer models.Offer
	if err := db.First(&offer, instance.OfferID).Error; err == nil {
		var nameValue string
		if input.Properties != nil {
			// Supprime toutes les propriétés existantes
			db.Where("instance_id = ?", instance.ID).Delete(&models.InstanceProperty{})
			// Ajoute les nouvelles propriétés
			for _, prop := range input.Properties {
				instProp := models.InstanceProperty{
					InstanceID:      instance.ID,
					OfferPropertyID: prop.OfferPropertyID,
					Value:           prop.Value,
				}
				db.Create(&instProp)
				if offer.NamePropertyID != nil && prop.OfferPropertyID == *offer.NamePropertyID {
					nameValue = prop.Value
				}
			}
			// Met à jour le nom si trouvé
			if nameValue != "" {
				instance.Name = nameValue
				db.Save(&instance)
			}
		}
	} else if input.Properties != nil {
		// fallback: juste mettre à jour les propriétés
		db.Where("instance_id = ?", instance.ID).Delete(&models.InstanceProperty{})
		for _, prop := range input.Properties {
			instProp := models.InstanceProperty{
				InstanceID:      instance.ID,
				OfferPropertyID: prop.OfferPropertyID,
				Value:           prop.Value,
			}
			db.Create(&instProp)
		}
	}
	db.Preload("Properties").First(&instance, instance.ID)
	c.JSON(http.StatusOK, instance)
}

// Suppression d'une instance
func DeleteInstance(c *gin.Context) {
	db := db.GetDB()
	var instance models.Instance
	id := c.Param("id")
	if err := db.First(&instance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}
	if err := db.Delete(&instance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
