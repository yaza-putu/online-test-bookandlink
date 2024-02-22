package routes

import (
	"github.com/labstack/echo/v4"
	queue "github.com/yaza-putu/online-test-bookandlink/internal/app/queue/handler"
)

var queueHandler = queue.NewQueueHandler()

func Api(r *echo.Echo) {
	route := r.Group("api")
	{
		v1 := route.Group("/v1")
		{
			v1.POST("/queue", queueHandler.Create)
			v1.GET("/queue/check", queueHandler.Recheck)
			v1.GET("/queue/rollback", queueHandler.Rollback)

			v1.GET("/jobs/done", queueHandler.Done)
			v1.GET("/jobs/pending", queueHandler.Pending)
			v1.GET("/jobs/failed", queueHandler.Failed)
		}
	}
}
