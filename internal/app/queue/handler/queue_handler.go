package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/service"
	"github.com/yaza-putu/online-test-bookandlink/internal/http/response"
	"net/http"
	"sync"
)

type queueHandler struct {
	queue service.Queue
}

func NewQueueHandler() *queueHandler {
	return &queueHandler{
		queue: service.NewQueue(service.SetWorker(20)), // set number of worker
	}
}

func (q *queueHandler) Create(ctx echo.Context) error {
	jobs := make(chan string)
	wg := new(sync.WaitGroup)

	for i := 0; i < 100; i++ {
		email := fmt.Sprintf("user%d@example.com", i+1)
		go q.queue.Add(context.Background(), jobs, wg, email)
	}

	go q.queue.DispatchWorkers(jobs, wg)

	wg.Wait()

	return ctx.JSON(http.StatusOK, response.Api(
		response.SetCode(http.StatusOK),
		response.SetMessage("Send to queue successfully")),
	)
}
