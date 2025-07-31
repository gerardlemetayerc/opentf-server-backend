package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
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
		// Vérifie si la table users existe avant migration
		var usersTableExists bool
		db.Raw("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='users';").Scan(&usersTableExists)

		// Vérifie si la table groups existe avant migration
		var groupsTableExists bool
		db.Raw("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='groups';").Scan(&groupsTableExists)

		db.AutoMigrate(&models.Module{}, &models.Property{}, &models.Instance{}, &models.Job{}, &iam.OIDCConfig{}, &iam.AuthMethod{}, &models.Group{}, &models.User{})
		db.AutoMigrate(&models.User{}, &models.Group{}) // pour la table de jointure user_groups
		db.AutoMigrate(&models.Domain{}, &models.SuggestedValue{})
		db.AutoMigrate(&models.OfferCategory{}, &models.Offer{})
		db.AutoMigrate(&models.OfferProperty{})
		db.AutoMigrate(&models.Instance{}, &models.InstanceProperty{})

		// Injecte l'utilisateur localadmin@local uniquement si la table users vient d'être créée
		var localAdmin models.User
		if !usersTableExists {
			hash, err := bcrypt.GenerateFromPassword([]byte("local"), bcrypt.DefaultCost)
			if err == nil {
				localAdmin = models.User{
					FirstName:    "Local",
					LastName:     "Admin",
					Email:        "localadmin@local",
					Status:       "active",
					AuthSource:   "locale",
					PasswordHash: string(hash),
				}
				db.Create(&localAdmin)
			}
		} else {
			db.Where("email = ? AND auth_source = ?", "localadmin@local", "locale").First(&localAdmin)
		}

		// Injecte les rôles et le groupe administrateur si la table groups vient d'être créée
		if !groupsTableExists {
			// Crée le groupe administrateur avec le rôle administrateur
			adminGroup := models.Group{
				Name:  "administrateur",
				Roles: models.StringArray{"administrateur"},
			}
			db.Create(&adminGroup)
			// Ajoute l'utilisateur localadmin au groupe administrateur
			if localAdmin.ID != 0 {
				db.Model(&adminGroup).Association("Users").Append(&localAdmin)
			}
			// Crée un groupe utilisateurs avec le rôle utilisateur
			userGroup := models.Group{
				Name:  "utilisateur",
				Roles: models.StringArray{"utilisateur"},
			}
			db.Create(&userGroup)
		}

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
