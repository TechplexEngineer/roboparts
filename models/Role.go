package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Users       []User       `json:"users"       gorm:"many2many:roles_users;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:roles_permissions;"`
}

func (o Role) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}
