package entity

import "time"

const (
	PENDING    = "PENDING"
	PROCESSING = "PROCESSING"
)

type Job struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string    `json:"name" gorm:"name"`
	Payload   string    `gorm:"type:longtext" json:"payload"`
	Attempts  int       `gorm:"type:integer;unsigned" json:"attempts"`
	Status    string    `gorm:"type:enum('PENDING', 'PROCESSING')" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Jobs []Job
