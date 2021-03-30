package models

import (
	"encoding/json"
)

// Vendor is a company that sells COTS parts
type Vendor struct {
	Common
	Name       string     `json:"name"`
	PartPrefix string     `json:"part_prefix"`
	Notes      string     `json:"notes"`
	Parts      []COTSPart `json:"parts"`
	Orders     []Order    `json:"orders"`
}

func (o Vendor) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}
