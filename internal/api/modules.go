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
