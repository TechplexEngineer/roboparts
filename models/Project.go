package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	Common
	Name       string   `json:"name"`
	PartPrefix string   `json:"part_prefix"`
	Archived   bool     `json:"archived"`
	Notes      string   `json:"notes"`
	Parts      []Part   `json:"parts"`
	Orders     []*Order `json:"orders" gorm:"many2many:projects_orders;"`
}

func (o Project) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

func (o *Project) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
