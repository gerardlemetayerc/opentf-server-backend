package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"opentf-server/internal/models"
	"opentf-server/internal/models/iam"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open("sqlite3", "opentf.db")
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
		db.AutoMigrate(&models.Module{}, &models.Property{}, &models.Instance{}, &models.Job{}, &iam.OIDCConfig{}, &iam.AuthMethod{}, &models.Group{}, &models.User{})
		db.AutoMigrate(&models.User{}, &models.Group{}) // pour la table de jointure user_groups
		db.AutoMigrate(&models.Domain{}, &models.SuggestedValue{})
		db.AutoMigrate(&models.OfferCategory{}, &models.Offer{})
		db.AutoMigrate(&models.OfferProperty{})
		db.AutoMigrate(&models.Instance{}, &models.InstanceProperty{})

		// Seed catégories d'offre par défaut
		var countCat int
		db.Model(&models.OfferCategory{}).Count(&countCat)
		if countCat == 0 {
			db.Create(&models.OfferCategory{Name: "Infrastructure"})
			db.Create(&models.OfferCategory{Name: "Identité & accès"})
		}
		// Ajoute l'index unique sur (email, auth_source)
		db.Model(&models.User{}).AddUniqueIndex("idx_user_email_source", "email", "auth_source")
		// Seed OIDC method if not exists
		var count int
		db.Model(&iam.AuthMethod{}).Where("method = ?", "oidc").Count(&count)
		if count == 0 {
			db.Create(&iam.AuthMethod{Method: "oidc", Enabled: false})
		}
	})
	return db
}
