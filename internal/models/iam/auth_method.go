package iam

type AuthMethod struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Method  string `gorm:"unique;not null" json:"method"`
	Enabled bool   `json:"enabled"`
}
