package repository

import (
	"context"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/database"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/db"
)

type (
	DoneJob interface {
		Create(ctx context.Context, job entity.DoneJob) (entity.DoneJob, error)
		All(ctx context.Context, page int, take int) (db.Pagination, error)
	}
	doneJobRepository struct{}
)

func NewDoneJob() *doneJobRepository {
	return &doneJobRepository{}
}

// Create done job
func (d *doneJobRepository) Create(ctx context.Context, job entity.DoneJob) (entity.DoneJob, error) {
	r := database.Instance.WithContext(ctx).Create(&job)

	return job, r.Error
}

// All done jobs
func (d *doneJobRepository) All(ctx context.Context, page int, take int) (db.Pagination, error) {
	e := entity.DoneJobs{}

	var pagination db.Pagination
	var totalRow int64

	pagination.SetSort("done_jobs.created_at desc")

	r := database.Instance.WithContext(ctx).Model(&e)
	r.Scopes(pagination.Paginate(page, take))

	r.Count(&totalRow).Find(&e)

	pagination.Rows = e
	pagination.CalculatePage(float64(totalRow))

	return pagination, r.Error
}
