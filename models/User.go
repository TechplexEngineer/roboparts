package models

import (
	"encoding/json"
)

type User struct {
	Common
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"password_hash"`
	Roles        []Role       `json:"roles"       gorm:"many2many:roles_users;"`
	Permissions  []Permission `json:"permissions" gorm:"many2many:permissions_users;"`
}

func (o User) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}
