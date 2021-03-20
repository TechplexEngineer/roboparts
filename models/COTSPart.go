package models

import (
	"gorm.io/gorm"
)

// COTSPart is represents a single off the shelf part
type COTSPart struct {
	gorm.Model
	Name       string `json:"name"`
	PartNumber string `json:"part_number"`
	QtyPerUnit int    `json:"qty_per_unit"`
	UnitCost   int64  `json:"unit_cost"`
	Link       string `json:"link"`
	Notes      string `json:"notes"`
	VendorID   uint   `json:"vendor_id"`
	Vendor     Vendor `json:"vendor"`
}

type COTSParts []COTSPart
