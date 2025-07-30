package models

import "time"

// Instance liée à une offre
// Statut par défaut : 'draft'
type Instance struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	OfferID     uint               `json:"offer_id"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Status      string             `json:"status" gorm:"default:draft"`
	RequesterID uint               `json:"requester_id"`
	ValidatorID uint               `json:"validator_id"`
	Properties  []InstanceProperty `gorm:"foreignKey:InstanceID" json:"properties"`
	Name        string             `json:"name"`
}

type InstanceProperty struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	InstanceID      uint   `json:"instance_id"`
	OfferPropertyID uint   `json:"offer_property_id"`
	Value           string `json:"value"`
}
