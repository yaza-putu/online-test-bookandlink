package service

import (
	"context"
	"fmt"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/repository"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/logger"
	"github.com/yaza-putu/online-test-bookandlink/pkg/unique"
	"math"
	"sync"
	"time"
)

type (
	Queue interface {
		Add(ctx context.Context, jobs chan<- []string, wg *sync.WaitGroup, emails []string) // add job to queue
		DispatchWorkers(jobs <-chan []string, wg *sync.WaitGroup)                           // dispatch worker
		sendEmail(workerIndex int, email string)                                            // send email job
	}
	queueService struct {
		totalWorker         int
		doneJobRepository   repository.DoneJob
		jobRepository       repository.Job
		failedJobRepository repository.FailedJob
	}
	optFunc func(*queueService)
)

func defaultOption() queueService {
	return queueService{
		totalWorker:         10,
		doneJobRepository:   repository.NewDoneJob(),
		jobRepository:       repository.NewJob(),
		failedJobRepository: repository.NewFailedJob(),
	}
}

func NewQueue(opts ...optFunc) *queueService {
	o := defaultOption()

	for _, fn := range opts {
		fn(&o)
	}

	return &queueService{
		totalWorker:         o.totalWorker,
		doneJobRepository:   repository.NewDoneJob(),
		jobRepository:       repository.NewJob(),
		failedJobRepository: repository.NewFailedJob(),
	}
}

func SetWorker(workers int) optFunc {
	return func(q *queueService) {
		q.totalWorker = workers
	}
}

// Mock repository for unit testing
func Mock(job repository.Job, doneJob repository.DoneJob, failedJob repository.FailedJob) optFunc {
	return func(q *queueService) {
		q.failedJobRepository = failedJob
		q.jobRepository = job
		q.doneJobRepository = doneJob
	}
}

// Add job to queue & send to worker
func (q *queueService) Add(ctx context.Context, jobs chan<- []string, wg *sync.WaitGroup, emails []string) {
	for _, email := range emails {
		_, err := q.jobRepository.Add(ctx, entity.Job{
			ID:      unique.Uid(13),
			Name:    "Send email to " + email,
			Payload: email,
			Status:  entity.PENDING,
		})
		// send error to central logger handler
		logger.New(err)
	}

	jobs <- emails
	wg.Add(1)
	close(jobs)
}

func (q *queueService) DispatchWorkers(jobs <-chan []string, wg *sync.WaitGroup) {
	for workerIndex := 1; workerIndex <= q.totalWorker; workerIndex++ {
		go func(workerIndex int, jobs <-chan []string, wg *sync.WaitGroup) {
			for job := range jobs {
				for _, v := range job {
					q.sendEmail(workerIndex, v)
				}
				wg.Done()
			}
		}(workerIndex, jobs, wg)
	}
}

func (q *queueService) sendEmail(workerIndex int, email string) {
	start := time.Now()
	ctx := context.Background()
	// send event websocket
	job, err := q.jobRepository.TakeOne(ctx, email)
	if err != nil {
		logger.New(err)
	} else {
		// send email
		// after success send email
		err = q.jobRepository.Delete(ctx, job.ID)
		logger.New(err)

		// mark job to done
		duration := time.Since(start)
		_, err := q.doneJobRepository.Create(ctx, entity.DoneJob{
			ID:          unique.Uid(13),
			Name:        job.Name,
			WorkerIndex: workerIndex,
			Duration:    fmt.Sprintf("%d ms", int(math.Ceil(float64(duration.Milliseconds())))),
		})
		logger.New(err)
	}
}
