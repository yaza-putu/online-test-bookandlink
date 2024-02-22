package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/service"
	"github.com/yaza-putu/online-test-bookandlink/internal/http/request"
	"github.com/yaza-putu/online-test-bookandlink/internal/http/response"
	"net/http"
)

type (
	queueHandler struct {
		queue service.Queue
	}
	requestValidation struct {
		TotalJobs int `json:"total_jobs" form:"total_jobs" validate:"required"`
	}
	paginationvalidation struct {
		Page int `json:"page" query:"page"`
		Take int `json:"take" query:"take"`
	}
)

func NewQueueHandler() *queueHandler {
	return &queueHandler{
		queue: service.NewQueue(service.SetMaxWorker(20)), // set number of worker
	}
}

// Create job && run queue
func (q *queueHandler) Create(ctx echo.Context) error {
	r := requestValidation{}
	// check input type data
	err := ctx.Bind(&r)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.BadRequest(err))
	}

	// validation
	result, err := request.Validation(&r)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, result)
	}

	// run queue
	q.queue.Run()

	var emails []string
	// generate job
	for i := 0; i < r.TotalJobs; i++ {
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

	return ctx.JSON(http.StatusOK, response.Api(
		response.SetCode(http.StatusOK),
		response.SetMessage("Recheck job successfully"),
	))
}

// Rollback failed job to queue
func (q *queueHandler) Rollback(ctx echo.Context) error {
	// cancel process if user close connection
	err := q.queue.Rollback(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.Api(
			response.SetCode(500),
			response.SetError(err),
			response.SetMessage("Internal server error"),
		))
	}

	return ctx.JSON(http.StatusOK, response.Api(
		response.SetCode(http.StatusOK),
		response.SetMessage("Rollback all failed job to queue successfully"),
	))
}

// Done jobs
func (q *queueHandler) Done(ctx echo.Context) error {
	r := paginationvalidation{}
	err := ctx.Bind(&r)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.BadRequest(errors.New("Bad requests")))
	}

	// validation
	d, err := request.Validation(&r)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, d)
	}

	res := q.queue.AllDoneJob(ctx.Request().Context(), r.Page, r.Take)

	return ctx.JSON(res.Code, res)
}

// Failed jobs
func (q *queueHandler) Failed(ctx echo.Context) error {
	r := paginationvalidation{}
	err := ctx.Bind(&r)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.BadRequest(errors.New("Bad requests")))
	}

	// validation
	d, err := request.Validation(&r)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, d)
	}

	res := q.queue.AllFailedJob(ctx.Request().Context(), r.Page, r.Take)

	return ctx.JSON(res.Code, res)
}

// Pending jobs
func (q *queueHandler) Pending(ctx echo.Context) error {
	r := paginationvalidation{}
	err := ctx.Bind(&r)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.BadRequest(errors.New("Bad requests")))
	}

	// validation
	d, err := request.Validation(&r)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, d)
	}

	res := q.queue.AllPendingJob(ctx.Request().Context(), r.Page, r.Take)

	return ctx.JSON(res.Code, res)
}
