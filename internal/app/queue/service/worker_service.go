package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/repository"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/logger"
	"github.com/yaza-putu/online-test-bookandlink/pkg/unique"
	"gorm.io/gorm"
	"math"
	"time"
)

type workerService struct {
	ID                  int
	WorkerQueue         chan chan Job
	JobQueue            chan Job
	QuitChan            chan bool
	jobRepository       repository.Job
	failedJobRepository repository.FailedJob
	doneJobRepository   repository.DoneJob
}

func NewWorker(id int, workerQueue chan chan Job) *workerService {
	return &workerService{
		ID:                  id,
		WorkerQueue:         workerQueue,
		JobQueue:            make(chan Job),
		QuitChan:            make(chan bool),
		jobRepository:       repository.NewJob(),
		failedJobRepository: repository.NewFailedJob(),
		doneJobRepository:   repository.NewDoneJob(),
	}
}

// Start worker
func (w workerService) Start() {
	go func() {
		for {
			// Register the current worker to the worker pool
			w.WorkerQueue <- w.JobQueue

			select {
			case job := <-w.JobQueue:
				start := time.Now()
				ctx, cancel := context.WithCancel(context.Background())

				// working
				j, err := w.jobRepository.TakeOne(context.Background(), job.Email)
				logger.New(err)

				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					// we assume if failed send email add attempt counter
					if j.Attempts < 3 {
						j.Attempts += 1
						er := w.jobRepository.Update(context.Background(), j.ID, j)
						logger.New(er)
						w.Stop()
					} else {
						// send to failed job
						er := w.jobRepository.Delete(ctx, j.ID)
						if err != nil {
							logger.New(er)
							cancel()
						}

						_, er = w.failedJobRepository.Create(ctx, entity.FailedJob{
							ID:        unique.Uid(13),
							Name:      j.Name,
							Payload:   job.Email,
							Exception: err.Error(),
						})

						if er != nil {
							logger.New(er)
							w.Stop()
						}

					}
				}

				// we assume we have sent the email
				err = w.jobRepository.Delete(ctx, j.ID)
				if err != nil {
					cancel()
				}

				duration := time.Since(start)
				done := int(math.Ceil(float64(duration.Milliseconds())))

				_, err = w.doneJobRepository.Create(ctx, entity.DoneJob{
					ID:          unique.Uid(13),
					Name:        j.Name,
					Duration:    fmt.Sprintf("%d ms", done),
					WorkerIndex: w.ID,
				})

				if err != nil {
					w.Stop()
				}

				fmt.Printf("Worker %d send email to: %s done in %d ms\n", w.ID, job.Email, done)
			case <-w.QuitChan:
				// Exits a goroutine when it receives a signal to stop
				return
			}
		}
	}()
}

// Stop worker
func (w workerService) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
