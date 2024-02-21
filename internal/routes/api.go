package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/auth/handler"
	queue "github.com/yaza-putu/online-test-bookandlink/internal/app/queue/handler"
)

var authhandler = handler.NewAuthHandler()
var queueHandler = queue.NewQueueHandler()

func Api(r *echo.Echo) {
	route := r.Group("api")
	{
		v1 := route.Group("/v1")
		{
			v1.POST("/token", authhandler.Create)
			v1.PUT("/token", authhandler.Refresh)

			v1.GET("/queue", queueHandler.Create)
			v1.GET("/queue/re-check", queueHandler.Recheck)
		}
	}
}
