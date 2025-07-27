package models

import (
	"time"
)

type User struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `gorm:"not null" json:"email"`
	Status       string     `json:"status"`
	LastLogin    *time.Time `json:"last_login"`
	AuthSource   string     `gorm:"not null" json:"auth_source"` // locale, oidc, etc.
	PasswordHash string     `json:"-"`                           // Ne jamais exposer le hash dans l'API
	Groups       []Group    `gorm:"many2many:user_groups;" json:"groups"`
}

// Ajoute une contrainte d'unicité sur (email, auth_source) dans la migration.
