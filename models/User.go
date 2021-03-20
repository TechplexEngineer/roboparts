package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                string `json:"email"`
	PasswordHash         string `json:"password_hash"`
	Password             string `json:"-" gorm:"-"`
	PasswordConfirmation string `json:"-" gorm:"-"`
}

func (o User) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

type Users []User
