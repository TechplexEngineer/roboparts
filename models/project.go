package models

import (
	"encoding/json"
)

type Project struct {
	Common
	Name       string   `json:"name"        form:"Name"       validate:"required"   gorm:"uniqueIndex"`
	PartPrefix string   `json:"part_prefix" form:"PartPrefix" validate:"required"`
	Archived   bool     `json:"archived"    form:"Archived"`
	Notes      string   `json:"notes"       form:"Notes"      ui:"textarea"`
	Parts      []Part   `json:"parts"`
	Orders     []*Order `json:"orders" gorm:"many2many:projects_orders;"`
}

func (o Project) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}
