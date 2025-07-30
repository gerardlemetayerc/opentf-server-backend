package models

import "time"

// StateLock pour la gestion des verrous sur les statefiles
// Un lock est lié à une instance et empêche les accès concurrents
// Le champ Info peut contenir des infos sur le worker, le user, etc.
type StateLock struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	InstanceID uint      `json:"instance_id"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Info       string    `json:"info"`
}
