package models

import "encoding/json"

// COTSPart is represents a single off the shelf part
type COTSPart struct {
	Common
	Name       string `json:"name"`
	PartNumber string `json:"part_number"`
	QtyPerUnit int    `json:"qty_per_unit"`
	UnitCost   int64  `json:"unit_cost"`
	Link       string `json:"link"`
	Notes      string `json:"notes"`
	VendorID   uint   `json:"vendor_id"`
	Vendor     Vendor `json:"vendor"`
}

func (o COTSPart) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}
