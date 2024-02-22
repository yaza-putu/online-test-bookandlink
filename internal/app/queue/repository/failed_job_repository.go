package repository

import (
	"context"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/database"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/db"
	"github.com/yaza-putu/online-test-bookandlink/pkg/unique"
	"gorm.io/gorm"
)

type (
	FailedJob interface {
		Create(ctx context.Context, job entity.FailedJob) (entity.FailedJob, error)
		Rollback(ctx context.Context) error
		All(ctx context.Context, page int, take int) (db.Pagination, error)
	}
	failedJobRepository struct {
		db *gorm.DB
	}
)

func NewFailedJob() *failedJobRepository {
	return &failedJobRepository{}
}

// Create failed job
func (f *failedJobRepository) Create(ctx context.Context, fJob entity.FailedJob) (entity.FailedJob, error) {
	r := database.Instance.WithContext(ctx).Create(&fJob)

	return fJob, r.Error
}

// Rollback failed job to processing
func (f *failedJobRepository) Rollback(ctx context.Context) error {
	fjob := entity.FailedJobs{}
	r := database.Instance.WithContext(ctx).Order("created_at asc").Find(&fjob)
	if r.Error != nil {
		return r.Error
	}

	database.Instance.Transaction(func(tx *gorm.DB) error {
		jobs := entity.Jobs{}
		for _, job := range fjob {
			jobs = append(jobs, entity.Job{
				ID:      unique.Uid(13),
				Name:    "Re-send email to " + job.Payload,
				Payload: job.Payload,
				Status:  entity.PENDING,
			})
		}
		t := tx.Create(&jobs)
		if t.Error != nil {
			return t.Error
		}

		d := tx.Delete(&fjob)
		if d.Error != nil {
			return d.Error
		}

		return nil
	})

	return nil
}

// All failed job
func (f *failedJobRepository) All(ctx context.Context, page int, take int) (db.Pagination, error) {
	e := entity.FailedJobs{}

	var pagination db.Pagination
	var totalRow int64

	pagination.SetSort("failed_jobs.created_at asc")

	r := database.Instance.WithContext(ctx).Model(&e)
	r.Scopes(pagination.Paginate(page, take))

	r.Count(&totalRow).Find(&e)

	pagination.Rows = e
	pagination.CalculatePage(float64(totalRow))

	return pagination, r.Error
}
