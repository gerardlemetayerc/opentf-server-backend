package iam

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models/iam"
)

// GET /api/iam/auth_methods
func GetAuthMethods(c *gin.Context) {
	db := db.GetDB()
	var methods []iam.AuthMethod
	db.Find(&methods)
	c.JSON(http.StatusOK, methods)
}

// POST /api/iam/auth_methods
func SetAuthMethod(c *gin.Context) {
	db := db.GetDB()
	var input iam.AuthMethod
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var method iam.AuthMethod
	if err := db.Where("method = ?", input.Method).First(&method).Error; err != nil {
		// Not found, create
		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
		return
	}
	// Found, update
	method.Enabled = input.Enabled
	if err := db.Save(&method).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, method)
}
