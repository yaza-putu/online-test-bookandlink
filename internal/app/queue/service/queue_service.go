package service

import (
	"context"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/repository"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/logger"
	"github.com/yaza-putu/online-test-bookandlink/pkg/unique"
)

type (
	Queue interface {
		Run()               // run queue
		dispatch()          // send job to worker
		EnqueueJob(job Job) // add job
		Stop()              // stop queue
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
		jobRepository: repository.NewJob(),
	}
}

func SetMaxWorker(workers int) optFunc {
	return func(q *queueService) {
		q.MaxWorkers = workers
	}
}

// Mock repository for unit testing
func Mock(job repository.Job, doneJob repository.DoneJob, failedJob repository.FailedJob) optFunc {
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
