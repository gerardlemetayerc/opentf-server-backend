package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Group struct {
	ID    uint        `gorm:"primary_key" json:"id"`
	Name  string      `gorm:"unique;not null" json:"name"`
	Roles StringArray `gorm:"type:text" json:"roles"`
	Users []User      `gorm:"many2many:user_groups;" json:"-"`
}

// StringArray permet de stocker un []string en base sous forme de JSON
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, a)
}
