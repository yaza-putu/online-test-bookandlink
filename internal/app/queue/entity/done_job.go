package entity

import "time"

type DoneJob struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name        string    `json:"name" gorm:"name"`
	Duration    string    `json:"duration" gorm:"duration"`
	WorkerIndex int       `json:"worker_index" gorm:"worker_index"`
	CreatedAt   time.Time `json:"created_at"`
}

type DoneJobs []DoneJob
