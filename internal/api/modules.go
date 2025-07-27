package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

func ListModules(c *gin.Context) {
	db := db.GetDB()
	var modules []models.Module
	db.Find(&modules)
	c.JSON(http.StatusOK, modules)
}

func CreateModule(c *gin.Context) {
	db := db.GetDB()
	var module models.Module
	if err := c.ShouldBindJSON(&module); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&module)
	c.JSON(http.StatusCreated, module)
}
