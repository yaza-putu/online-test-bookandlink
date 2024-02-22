package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/entity"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/repository"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/logger"
	"gorm.io/gorm"
	"math"
	"time"
)

type workerService struct {
	ID            int
	WorkerQueue   chan chan Job
	JobQueue      chan Job
	QuitChan      chan bool
	jobRepository repository.Job
}

func NewWorker(id int, workerQueue chan chan Job) *workerService {
	return &workerService{
		ID:            id,
		WorkerQueue:   workerQueue,
		JobQueue:      make(chan Job),
		QuitChan:      make(chan bool),
		jobRepository: repository.NewJob(),
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
						j.Status = entity.FAILED
						er := w.jobRepository.Update(ctx, j.ID, j)

						if er != nil {
							logger.New(er)
							w.Stop()
							cancel()
						}

					}
				}

				// we assume we have sent the email
				// and add time sleep to make behavior real handling job
				time.Sleep(time.Millisecond * 700)
				// -------------------------------
				duration := time.Since(start)
				done := int(math.Ceil(float64(duration.Milliseconds())))

				j.Status = entity.DONE
				j.Duration = fmt.Sprintf("%d ms", done)
				j.WorkerIndex = w.ID
				err = w.jobRepository.Update(ctx, j.ID, j)
				if err != nil {
					logger.New(err)
					cancel()
				}

				if err != nil {
					w.Stop()
					logger.New(err)
				}

				fmt.Printf("Worker %d send email to: %s done in %d ms\n", w.ID, job.Email, done)
				if len(connections) > 0 {
					broadcastMessage("monitor", fmt.Sprintf("Worker %d send email to: %s done in %d ms\n", w.ID, job.Email, done))
					eventCountJob()
					eventUpdateTable(1, 10, "")
				}
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
