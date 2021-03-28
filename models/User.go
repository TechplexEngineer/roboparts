package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username             string       `json:"username"`
	Email                string       `json:"email"`
	PasswordHash         string       `json:"password_hash"`
	Password             string       `json:"-"           gorm:"-"`
	PasswordConfirmation string       `json:"-"           gorm:"-"`
	Roles                []Role       `json:"roles"       gorm:"many2many:roles_users;"`
	Permissions          []Permission `json:"permissions" gorm:"many2many:permissions_users;"`
}

func (o User) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}
