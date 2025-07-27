package models

type Domain struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"unique;not null" json:"name"`
	Label string `json:"label"`
}

type SuggestedValue struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	DomainID       uint    `json:"domain_id"`
	Domain         Domain  `gorm:"foreignKey:DomainID" json:"domain"`
	DisplayValue   string  `json:"display_value"`
	RealValue      string  `json:"real_value"`
	ParentDomainID *uint   `json:"parent_domain_id,omitempty"`
	ParentValue    *string `json:"parent_value,omitempty"`
}
