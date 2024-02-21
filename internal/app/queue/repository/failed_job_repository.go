package repository

import (
	"context"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/database"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/db"
	"gorm.io/gorm"
)

type (
	FailedJob interface {
		Create(ctx context.Context, job entity.FailedJob) (entity.FailedJob, error)
		Rollback(ctx context.Context, id string) error
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
func (f *failedJobRepository) Rollback(ctx context.Context, id string) error {
	r := database.Instance.WithContext(ctx).Where("id = ?", id).Delete(&entity.FailedJob{})

	return r.Error
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
