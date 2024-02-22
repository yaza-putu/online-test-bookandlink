package entity

import "time"

const (
	PENDING    = "PENDING"
	PROCESSING = "PROCESSING"
	DONE       = "DONE"
	FAILED     = "FAILED"
)

type Job struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name        string    `json:"name" gorm:"name"`
	Payload     string    `gorm:"type:longtext" json:"payload"`
	Attempts    int       `gorm:"type:integer;unsigned" json:"attempts"`
	Status      string    `gorm:"type:enum('PENDING', 'PROCESSING','DONE', 'FAILED')" json:"status"`
	Duration    string    `json:"duration" gorm:"duration"`
	WorkerIndex int       `json:"worker_index" gorm:"worker_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Jobs []Job
