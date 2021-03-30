package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	Common
	Status     string      `json:"status"`
	OrderedAt  *time.Time  `json:"ordered_at"`
	PaidForBy  string      `json:"paid_for_by"`
	TaxCost    int64       `json:"tax_cost"`
	Notes      string      `json:"notes"`
	OrderItems []OrderItem `json:"order_items"`
	VendorId   *uuid.UUID  `json:"vendorId"`
	Vendor     *Vendor     `json:"vendor"`
	Projects   []*Project  `json:"projects" gorm:"many2many:projects_orders;"`

	//	# The list of possible order statuses. Key: string stored in database, value: what is displayed to the user.
	//	STATUS_MAP = {
	//	"open" => "Open",
	//	"ordered" => "Ordered",
	//	"received" => "Received"
	//}

	//	def subtotal
	//	order_items.map(&:total_cost).inject(0) { |sum, cost| sum + cost }
	//	end
	//
	//	def total_cost
	//	subtotal + tax_cost.to_f + shipping_cost.to_f
	//	end
}

func (o Order) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
