package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

// Represents a line item in an order from a vendor.
type OrderItem struct {
	gorm.Model
	Quantity    int64   `json:"quantity"`
	Description string  `json:"description"`
	UnitCost    int64   `json:"unitCost"`
	Notes       string  `json:"notes"`
	ProjectID   uint    `json:"projectId"`
	Project     Project `json:"project"`
	OrderID     uint    `json:"orderId"`
	Order       Order   `json:"order"`
	PartID      uint    `json:"partId"`
	Part        Part    `json:"part"`
}

func (o OrderItem) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

func (o OrderItem) TotalCost() int64 {
	return o.UnitCost * o.Quantity
}
