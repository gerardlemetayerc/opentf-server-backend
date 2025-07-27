package models

type Property struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	ModuleID    uint   `json:"module_id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // string, number, bool, etc.
	Required    bool   `json:"required"`
	Default     string `json:"default"`
	Description string `json:"description"`
}
