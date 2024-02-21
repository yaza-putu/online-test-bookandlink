package entity

import "time"

type FailedJob struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string    `json:"name" gorm:"name"`
	Payload   string    `gorm:"type:longtext" json:"payload"`
	Exception string    `gorm:"type:longtext" json:"exception"`
	CreatedAt time.Time `json:"created_at"`
}

type FailedJobs []FailedJob
