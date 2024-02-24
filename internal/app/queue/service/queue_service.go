package service

import (
	"context"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/repository"
	"github.com/yaza-putu/online-test-bookandlink/internal/http/response"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/logger"
	"github.com/yaza-putu/online-test-bookandlink/pkg/unique"
	"net/http"
)

type (
	Queue interface {
		Run()                                                                            // run queue
		dispatch()                                                                       // send job to worker
		EnqueueJob(job Job)                                                              // add job
		Stop()                                                                           // stop queue
		Check()                                                                          // check queue pending
		Rollback(ctx context.Context) error                                              // rollback all failed job to queue
		AllJob(ctx context.Context, page int, take int, keyword string) response.DataApi // all done
	}
	Job struct {
		Email string
	}
	queueService struct {
		WorkerPool    chan chan Job
		MaxWorkers    int
		JobQueue      chan Job
		Quit          chan bool
		jobRepository repository.Job
	}
	optFunc func(*queueService)
)

func defaultOption() queueService {
	return queueService{
		MaxWorkers:    10,
		WorkerPool:    make(chan chan Job, 10),
		JobQueue:      make(chan Job),
		Quit:          make(chan bool),
		jobRepository: repository.NewJob(),
	}
}

func NewQueue(opts ...optFunc) *queueService {
	o := defaultOption()

	for _, fn := range opts {
		fn(&o)
	}

	return &queueService{
		MaxWorkers:    o.MaxWorkers,
		WorkerPool:    o.WorkerPool,
		JobQueue:      o.JobQueue,
		Quit:          o.Quit,
		jobRepository: o.jobRepository,
	}
}

func SetMaxWorker(workers int) optFunc {
	return func(q *queueService) {
		q.MaxWorkers = workers
		q.WorkerPool = make(chan chan Job, workers)
	}
}

// Mock repository for unit testing
func Mock(job repository.Job) optFunc {
	return func(q *queueService) {
		q.jobRepository = job
	}
}

// Run queue
func (q *queueService) Run() {
	for i := 0; i < q.MaxWorkers; i++ {
		worker := NewWorker(i, q.WorkerPool)
		worker.Start()
	}

	go q.dispatch()
}

// dispatch job to worker
func (q *queueService) dispatch() {
	for {
		select {
		case job := <-q.JobQueue:
			go func(job Job) {
				workerJobQueue := <-q.WorkerPool
				workerJobQueue <- job
			}(job)
		case <-q.Quit:
			return
		}
	}
}

// EnqueueJob in the queue
func (q *queueService) EnqueueJob(job Job) {
	_, err := q.jobRepository.Add(context.Background(), entity.Job{
		ID:      unique.Uid(13),
		Name:    "Send email to " + job.Email,
		Payload: job.Email,
		Status:  entity.PENDING,
	})
	logger.New(err)

	// if no error send to job queue
	if err == nil {
		q.JobQueue <- job
	}

}

// Stop queue
func (q *queueService) Stop() {
	go func() {
		q.Quit <- true
	}()
}

// Check all job pending
func (q *queueService) Check() {
	j, err := q.jobRepository.Pending(context.Background())
	logger.New(err)

	for _, v := range j {
		q.JobQueue <- Job{Email: v.Payload}
	}
}

// Rollback failed job to queue
func (q *queueService) Rollback(ctx context.Context) error {
	err := q.jobRepository.Rollback(ctx)

	return err
}

// AllJobs
func (q *queueService) AllJob(ctx context.Context, page int, take int, keyword string) response.DataApi {
	r, err := q.jobRepository.All(ctx, page, take, keyword)
	if err != nil {
		return response.Api(
			response.SetCode(http.StatusInternalServerError),
			response.SetMessage("Internal server error"),
			response.SetError(err),
		)
	}

	return response.Api(
		response.SetCode(http.StatusOK),
		response.SetData(r),
	)
}
