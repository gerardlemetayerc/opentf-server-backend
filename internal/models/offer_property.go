package models

type OfferProperty struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	OfferID        uint   `gorm:"index" json:"offer_id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Label          string `json:"label"`
	Description    string `json:"description"`
	Required       bool   `json:"required"`
	DefaultValue   string `json:"default_value"`
	MinValue       string `json:"min_value"`
	MaxValue       string `json:"max_value"`
	MetadataSource string `json:"metadata_source"`
	DependsOn      string `json:"depends_on"`
	CustomJS       string `json:"customjs"`
}
