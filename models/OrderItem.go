package models

import (
	"encoding/json"
	"github.com/google/uuid"
)

// Represents a line item in an order from a vendor.
type OrderItem struct {
	Common
	Quantity    int64     `json:"quantity"`
	Description string    `json:"description"`
	UnitCost    int64     `json:"unitCost"`
	Notes       string    `json:"notes"`
	ProjectID   uuid.UUID `json:"projectId;type:uuid"`
	Project     Project   `json:"project"`
	OrderID     uuid.UUID `json:"orderId;type:uuid"`
	Order       Order     `json:"order"`
	PartID      uuid.UUID `json:"partId;type:uuid"`
	Part        Part      `json:"part"`
}

func (o OrderItem) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

func (o OrderItem) TotalCost() int64 {
	return o.UnitCost * o.Quantity
}
