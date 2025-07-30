package models

import (
	"time"
)

type Module struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Source      string     `json:"source"`         // URL or path to the module
	Cost        *float64   `json:"cost,omitempty"` // optional cost
	Properties  []Property `gorm:"foreignkey:ModuleID" json:"properties"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Job struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	InstanceID uint      `json:"instance_id"`
	Action     string    `json:"action"` // apply, destroy, update
	Status     string    `json:"status"` // queued, running, done, error
	Log        string    `json:"log"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
