package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type User struct {
	Common
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"password_hash"`
	Password     string       `json:"-"           gorm:"-"`
	Roles        []Role       `json:"roles"       gorm:"many2many:roles_users;"`
	Permissions  []Permission `json:"permissions" gorm:"many2many:permissions_users;"`
}

func (o User) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

func (o *User) BeforeCreate(tx *gorm.DB) error {
	if o.Password == "" {
		return nil // nothing to do
	}

	return nil
}
