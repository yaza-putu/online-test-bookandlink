package repository

import (
	"context"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/database"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/db"
	"gorm.io/gorm/clause"
)

type (
	Job interface {
		Add(ctx context.Context, job entity.Job) (entity.Job, error)
		TakeOne(ctx context.Context, email string) (entity.Job, error)
		Delete(ctx context.Context, id string) error
		All(ctx context.Context, page int, take int) (db.Pagination, error)
	}
	jobRepository struct{}
)

// NewJob instance
func NewJob() *jobRepository {
	return &jobRepository{}
}

// Add new job to queue
func (j *jobRepository) Add(ctx context.Context, job entity.Job) (entity.Job, error) {
	r := database.Instance.WithContext(ctx).Create(&job)
	return job, r.Error
}

// TakeOne to update status
func (j *jobRepository) TakeOne(ctx context.Context, email string) (entity.Job, error) {
	e := entity.Job{}
	r := database.Instance.Clauses(clause.Locking{Strength: "UPDATE"}).Where("payload = ?", email).Where("status = ?", entity.PENDING).First(&e)

	e.Status = entity.PROCESSING
	database.Instance.WithContext(ctx).Where("id = ?", e.ID).Updates(&e)
	return e, r.Error
}

// Delete processed job
func (j *jobRepository) Delete(ctx context.Context, id string) error {
	r := database.Instance.WithContext(ctx).Where("id = ?", id).Delete(&entity.Job{})

	return r.Error
}

// All job with pagination
func (j *jobRepository) All(ctx context.Context, page int, take int) (db.Pagination, error) {
	e := entity.Jobs{}

	var pagination db.Pagination
	var totalRow int64

	pagination.SetSort("jobs.created_at asc")

	r := database.Instance.WithContext(ctx).Model(&e)
	r.Scopes(pagination.Paginate(page, take))

	r.Count(&totalRow).Find(&e)

	pagination.Rows = e
	pagination.CalculatePage(float64(totalRow))

	return pagination, r.Error
}