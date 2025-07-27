package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
	"time"
)

// GET /api/users/:id
func GetUser(c *gin.Context) {
	db := db.GetDB()
	var user models.User
	id := c.Param("id")
	if err := db.Preload("Groups").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

// PUT /api/users/:id
func UpdateUser(c *gin.Context) {
	db := db.GetDB()
	var user models.User
	id := c.Param("id")
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if v, ok := body["first_name"].(string); ok {
		user.FirstName = v
	}
	if v, ok := body["last_name"].(string); ok {
		user.LastName = v
	}
	if v, ok := body["status"].(string); ok {
		user.Status = v
	}
	if v, ok := body["password"].(string); ok && user.AuthSource == "locale" {
		hash, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hash failed"})
			return
		}
		user.PasswordHash = string(hash)
	}
	// Gestion des groupes
	if groupsRaw, ok := body["groups"]; ok {
		var groupIDs []uint
		switch v := groupsRaw.(type) {
		case []interface{}:
			for _, g := range v {
				switch gVal := g.(type) {
				case map[string]interface{}:
					if idFloat, ok := gVal["id"].(float64); ok {
						groupIDs = append(groupIDs, uint(idFloat))
					}
				case float64:
					groupIDs = append(groupIDs, uint(gVal))
				case int:
					groupIDs = append(groupIDs, uint(gVal))
				}
			}
		}
		var groups []models.Group
		if len(groupIDs) > 0 {
			if err := db.Where("id IN (?)", groupIDs).Find(&groups).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ids"})
				return
			}
		}
		if err := db.Model(&user).Association("Groups").Replace(groups).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.PasswordHash = ""
	// Recharge les groupes pour la réponse
	db.Model(&user).Association("Groups").Find(&user.Groups)
	c.JSON(http.StatusOK, user)
}

// DELETE /api/users/:id
func DeleteUser(c *gin.Context) {
	db := db.GetDB()
	var user models.User
	id := c.Param("id")
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// GET /api/users
func GetUsers(c *gin.Context) {
	db := db.GetDB()
	var users []models.User
	db.Preload("Groups").Find(&users)
	for i := range users {
		users[i].PasswordHash = ""
	}
	c.JSON(http.StatusOK, users)
}

// POST /api/users (création)
func CreateUser(c *gin.Context) {
	db := db.GetDB()
	var input struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		Status     string `json:"status"`
		AuthSource string `json:"auth_source"`
		Password   string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		Email:      input.Email,
		Status:     input.Status,
		AuthSource: input.AuthSource,
	}
	if input.AuthSource == "locale" && input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hash failed"})
			return
		}
		user.PasswordHash = string(hash)
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

// POST /api/users/login (auth locale)
func LoginUser(c *gin.Context) {
	db := db.GetDB()
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := db.Where("email = ? AND auth_source = ?", input.Email, "locale").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	now := time.Now()
	db.Model(&user).Update("last_login", &now)
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}
