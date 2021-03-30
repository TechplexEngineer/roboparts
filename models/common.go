package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Common struct {
	ID        uuid.UUID `gorm:"primarykey;type:VARCHAR(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (o *Common) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
