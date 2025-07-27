package iam

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models/iam"
)

// GET /api/iam/auth/oidc
func GetOIDCConfig(c *gin.Context) {
	db := db.GetDB()
	var config iam.OIDCConfig
	if err := db.First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "OIDC config not found"})
		return
	}
	c.JSON(http.StatusOK, config)
}

// POST /api/iam/auth/oidc
func SetOIDCConfig(c *gin.Context) {
	db := db.GetDB()
	var input iam.OIDCConfig
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var config iam.OIDCConfig
	if err := db.First(&config).Error; err != nil {
		// Not found, create
		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
		return
	}
	// Found, update
	// On met à jour explicitement tous les champs, y compris Strict
	config.LocalEnabled = input.LocalEnabled
	config.OIDCEnabled = input.OIDCEnabled
	config.Issuer = input.Issuer
	config.ClientID = input.ClientID
	config.ClientSecret = input.ClientSecret
	config.Scopes = input.Scopes
	config.TokenEndpoint = input.TokenEndpoint
	config.AuthorizationEndpoint = input.AuthorizationEndpoint
	config.Strict = input.Strict
	config.Claims = input.Claims
	if err := db.Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}
