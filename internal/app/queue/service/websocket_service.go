package service

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/novalagung/gubrak/v2"
	"github.com/yaza-putu/online-test-bookandlink/internal/app/queue/repository"
	"github.com/yaza-putu/online-test-bookandlink/internal/pkg/logger"
	"log"
	"strings"
	"sync"
)

// event name
const (
	DASHBOARD = "dashboard_count"
	CLOSE     = "close"
	TABLE     = "update_table"
)

var (
	connections = make([]*WsConn, 0)
	mutex       sync.Mutex
	pageNum     int
)

type (
	WsConn struct {
		*websocket.Conn
	}
	SocketPayload struct {
		Event string `json:"event"`
		Data  any    `json:"data"`
	}
	SocketResponse struct {
		Event string `json:"event"`
		Data  any    `json:"data"`
	}
)

func NewWsService(ctx echo.Context, ws *websocket.Conn) {
	mutex.Lock()
	currentConn := WsConn{Conn: ws}
	connections = append(connections, &currentConn)
	mutex.Unlock()

	go handleIO(ctx, &currentConn, connections)
}

func handleIO(ctx echo.Context, currentConn *WsConn, connections []*WsConn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	payloadForMe := SocketPayload{}
	errForME := currentConn.ReadJSON(&payloadForMe)
	if errForME != nil {
		log.Println(errForME)
	}

	broadcastMessage("client_connect", len(connections))
	eventUpdateTable(1, 10, "")
	eventCountJob()

	manageEvent(currentConn, ctx, payloadForMe)

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				broadcastMessage(CLOSE, "")
				closeConnection(currentConn)
				return
			}

			log.Println("ERROR", err.Error())
			continue
		}

		manageEvent(currentConn, ctx, payload)
	}
}

func manageEvent(currentConn *WsConn, ctx echo.Context, wssReq SocketPayload) {
	if wssReq.Event == "pagination" {
		eventUpdateTable(int(wssReq.Data.(float64)), 10, "")
	}
}

func broadcastMessage(kind string, data any) {
	for _, eachConn := range connections {
		mutex.Lock()
		eachConn.WriteJSON(SocketResponse{
			Event: kind,
			Data:  data,
		})
		mutex.Unlock()
	}
}

func closeConnection(currentConn *WsConn) {
	filtered := gubrak.From(connections).Reject(func(each *WsConn) bool {
		return each == currentConn
	}).Result()
	connections = filtered.([]*WsConn)
}

// event count job
func eventCountJob() {
	jobRepository := repository.NewJob()
	p, d, f := jobRepository.CountJob()
	count := struct {
		Done    int
		Pending int
		Failed  int
	}{Done: d, Pending: p, Failed: f}

	broadcastMessage(DASHBOARD, count)
}

func eventUpdateTable(page int, take int, q string) {
	pageNum = page
	jobRepository := repository.NewJob()

	data, err := jobRepository.All(context.Background(), page, take, q)
	logger.New(err)
	if err == nil {
		broadcastMessage(TABLE, data)
	}
}

func eventWorkerLog(msg string) {
	broadcastMessage("worker_log", msg)
}
