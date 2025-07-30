package models

type OfferCategory struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}

type Offer struct {
	ID                uint          `gorm:"primaryKey" json:"id"`
	Name              string        `json:"name"`
	Version           string        `json:"version"`
	Icon              string        `json:"icon"`
	GitURL            string        `json:"git_url"`
	Active            bool          `json:"active"`
	CategoryID        uint          `json:"category_id"`
	Category          OfferCategory `gorm:"foreignKey:CategoryID" json:"category"`
	AutoValidated     bool          `json:"auto_validated"`
	ValidationGroupID *uint         `json:"validation_group_id"`
	NamePropertyID    *uint         `json:"name_property_id"`
	ModuleID          *uint         `json:"module_id"` // lien vers le module
}
