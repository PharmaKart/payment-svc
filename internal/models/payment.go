package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID       uuid.UUID `gorm:"not null;unique"`
	CustomerID    uuid.UUID `gorm:"not null"`
	TransactionID string    `gorm:"not null;unique"`
	Amount        float64   `gorm:"not null"`
	Status        string    `gorm:"type:varchar(50);not null;check:status IN ('pending', 'successful', 'failed')"`
	CreatedAt     time.Time `gorm:"default:now()"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
