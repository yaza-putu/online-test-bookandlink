package handler

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/service"
	"net/http"
)

type wssHandler struct{}

func NewWsHandler() *wssHandler {
	return &wssHandler{}
}

func (w *wssHandler) Connect(ctx echo.Context) error {
	ws, err := websocket.Upgrade(ctx.Response(), ctx.Request(), ctx.Response().Header(), 1024, 1024)
	if err != nil {
		http.Error(ctx.Response(), "Could not open websocket connection", http.StatusBadRequest)
	}
	if err != nil {
		return err
	}

	service.NewWsService(ctx, ws)
	return nil
}
