package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/service"
	"github.com/yaza-putu/online-test-bookandlink/internal/http/response"
	"net/http"
)

type queueHandler struct {
	queue service.Queue
}

func NewQueueHandler() *queueHandler {
	return &queueHandler{
		queue: service.NewQueue(service.SetMaxWorker(20)), // set number of worker
	}
}

// Create job && run queue
func (q *queueHandler) Create(ctx echo.Context) error {

	q.queue.Run()

	var emails []string

	// generate job
	for i := 0; i < 1000; i++ {
		emails = append(emails, fmt.Sprintf("user%d@example.com", i+1))
	}

	// add job
	go func(emails []string) {
		for _, email := range emails {
			q.queue.EnqueueJob(service.Job{Email: email})
		}
	}(emails)

	return ctx.JSON(http.StatusOK, response.Api(
		response.SetCode(http.StatusOK),
		response.SetMessage("Create job successfully")),
	)
}

// Recheck job pending
func (q *queueHandler) Recheck(ctx echo.Context) error {
	q.queue.Run()

	go q.queue.Check()

	return ctx.JSON(http.StatusOK, response.Api(response.SetCode(http.StatusOK), response.SetMessage("Recheck job successfully")))
}
