package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
)

func ListJobs(c *gin.Context) {
	db := db.GetDB()
	var jobs []models.Job
	db.Find(&jobs)
	c.JSON(http.StatusOK, jobs)
}

func CreateJob(c *gin.Context) {
	db := db.GetDB()
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&job)
	c.JSON(http.StatusCreated, job)
}
