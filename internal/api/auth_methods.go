package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

// POST /api/users/login_oidc
func LoginOIDCUser(c *gin.Context) {
	// Exemple : reçoit un id_token OIDC, le valide, extrait l'email, vérifie l'utilisateur
	var input struct {
		IDToken string `json:"id_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Ici, ajouter la logique de validation OIDC (non implémentée)
	// ...
	c.JSON(http.StatusNotImplemented, gin.H{"error": "OIDC login not implemented"})
}

// POST /api/users/login_token
func LoginTokenUser(c *gin.Context) {
	db := db.GetDB()
	var input struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := db.Where("api_token = ?", input.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}
