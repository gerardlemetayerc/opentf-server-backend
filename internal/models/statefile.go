package models

import "time"

// StateFile stocké en base
// Peut être lié à une instance, un module, etc.
type StateFile struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	InstanceID uint      `json:"instance_id"`
	Data       []byte    `json:"data" gorm:"type:blob"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
