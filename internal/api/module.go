package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
	"os"
	"os/exec"
	"path/filepath"
)

// POST /api/modules/:id/update
func UpdateModuleArchive(c *gin.Context) {
	db := db.GetDB()
	var module models.Module
	id := c.Param("id")
	if err := db.First(&module, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}
	if module.GitURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No GitURL for this module"})
		return
	}
	tmpDir, err := ioutil.TempDir("", "moduleclone")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp dir"})
		return
	}
	defer os.RemoveAll(tmpDir)
	cmd := exec.Command("git", "clone", module.GitURL, tmpDir)
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Git clone failed"})
		return
	}
	archivePath := filepath.Join(tmpDir, "module.zip")
	err = zipFolder(tmpDir, archivePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Zip failed"})
		return
	}
	archiveData, err := ioutil.ReadFile(archivePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Read zip failed"})
		return
	}
	module.Data = archiveData
	if err := db.Save(&module).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": true})
}

// GET /api/modules/:id/archive
func GetModuleArchive(c *gin.Context) {
	db := db.GetDB()
	var module models.Module
	id := c.Param("id")
	if err := db.First(&module, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}
	if len(module.Data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No archive available for this module"})
		return
	}
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=module.zip")
	c.Writer.Write(module.Data)
}

// GET /api/modules
func GetModules(c *gin.Context) {
	db := db.GetDB()
	var modules []models.Module
	db.Find(&modules)
	c.JSON(http.StatusOK, modules)
}

// POST /api/modules
func CreateModule(c *gin.Context) {
	db := db.GetDB()
	var input models.Module
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
