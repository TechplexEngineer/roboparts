package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name        string
	Description string
	Users       []User `json:"users" gorm:"many2many:permissions_users;"`
	Roles       []Role `json:"roles" gorm:"many2many:roles_permissions;"`
}

func (o Permission) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}
