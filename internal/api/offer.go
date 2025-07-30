package api

import (
	"archive/zip"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"opentf-server/internal/db"
	"opentf-server/internal/models"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	// Par défaut, une offre est auto-validée
	input.AutoValidated = true

	// Si un GitURL de module est précisé, on clone et stocke le contenu
	if input.GitURL != "" {
		tmpDir, err := ioutil.TempDir("", "moduleclone")
		if err == nil {
			cmd := exec.Command("git", "clone", input.GitURL, tmpDir)
			if err := cmd.Run(); err == nil {
				archivePath := filepath.Join(tmpDir, "module.zip")
				err := zipFolder(tmpDir, archivePath)
				if err == nil {
					archiveData, err := ioutil.ReadFile(archivePath)
					if err == nil {
						module := models.Module{
							Name:        input.Name,
							Description: "Module importé automatiquement",
							GitURL:      input.GitURL,
							Version:     input.Version,
							Active:      true,
							Data:        archiveData,
						}
						db.Create(&module)
						input.ModuleID = &module.ID
					}
				}
			}
			os.RemoveAll(tmpDir)
		}
	}

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Preload("Category").First(&input, input.ID)
	c.JSON(http.StatusOK, input)
}

// zipFolder zips the contents of srcDir into destZip
func zipFolder(srcDir, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore the zip file itself
		if path == destZip {
			return nil
		}
		relPath := strings.TrimPrefix(path, srcDir)
		relPath = strings.TrimLeft(relPath, string(filepath.Separator))
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		fHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		fHeader.Name = relPath
		fHeader.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(fHeader)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		return err
	})
	return err
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
	offer.AutoValidated = input.AutoValidated
	offer.ValidationGroupID = input.ValidationGroupID
	offer.NamePropertyID = input.NamePropertyID
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
