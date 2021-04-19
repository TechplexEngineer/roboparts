package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Common struct {
	ID        uuid.UUID      `gorm:"primarykey;type:VARCHAR(36)" ui:"-"`
	CreatedAt time.Time      `ui:"-"`
	UpdatedAt time.Time      `ui:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" ui:"-"`
}

func (o *Common) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
