// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package net

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zhouhp1295/g3"
	"github.com/zhouhp1295/g3/helpers"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

var workers map[string]bool

type WsConnStatus int

const (
	Connecting WsConnStatus = iota
	Connected
	Closed
)

type WsConn struct {
	Uuid     string
	Conn     *websocket.Conn
	Query    map[string]string
	AuthData map[string]interface{}
	CreateAt time.Time
	Status   WsConnStatus
}

type WsWorker struct {
	Router      string
	Upgrader    websocket.Upgrader
	OnConnected func(conn *WsConn)
	OnClosed    func(conn *WsConn)
	OnError     func(conn *WsConn, err error)
	OnAuth      func(conn *WsConn) (map[string]interface{}, bool)
	OnMessage   func(conn *WsConn, message []byte)
	ConnHandler WsConnHandler
	connections map[string]*WsConn
	rwMutex     sync.RWMutex
}

func (w *WsWorker) listen(conn *WsConn) {
	g3.ZL().Info("start listen", zap.String("uuid", conn.Uuid))
	for {
		_, message, err := conn.Conn.ReadMessage()
		if err != nil {
			w.OnError(conn, err)
			break
		}
		g3.ZL().Debug("on message", zap.String("uuid", conn.Uuid))
		if w.OnMessage != nil {
			w.OnMessage(conn, message)
		} else {
			_ = conn.Conn.WriteJSON(map[string]interface{}{
				"msg": "Message From Server : " + string(message),
			})
		}
	}
}

// handleConnect
func (w *WsWorker) handleConnect(conn *WsConn) bool {
	g3.ZL().Info("connected",
		zap.String("uuid", conn.Uuid),
		zap.Reflect("query", conn.Query),
	)
	if w.OnAuth != nil {
		var ok bool
		conn.AuthData, ok = w.OnAuth(conn)
		if !ok {
			return false
		}
	} else {
		conn.AuthData = make(map[string]interface{})
	}
	w.rwMutex.Lock()
	w.connections[conn.Uuid] = conn
	conn.Status = Connected
	w.rwMutex.Unlock()
	return true
}

func (w *WsWorker) closeConn(conn *WsConn) {
	g3.ZL().Info("connected",
		zap.String("uuid", conn.Uuid),
		zap.Reflect("query", conn.Query),
	)
	w.rwMutex.Lock()
	if Closed != conn.Status {
		if w.OnClosed != nil {
			w.OnClosed(conn)
		}
		if conn != nil {
			_ = conn.Conn.Close()
		}
		delete(w.connections, conn.Uuid)
	}
	w.rwMutex.Unlock()
}

// WsConnHandler WsConnHandler
type WsConnHandler func(conn *WsConn)

// WsWorkerOption WsWorkerOption
type WsWorkerOption func(*WsWorker)

// WithUpgrader WithUpgrader
func WithUpgrader(up websocket.Upgrader) WsWorkerOption {
	return func(worker *WsWorker) {
		worker.Upgrader = up
	}
}

// WithWsHandler WithWsHandler
func WithWsHandler(handler WsConnHandler) WsWorkerOption {
	return func(worker *WsWorker) {
		worker.ConnHandler = handler
	}
}

// defaultUpgrader defaultUpgrader
func defaultUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func init() {
	workers = make(map[string]bool)
}

func onError(w *WsWorker, conn *WsConn, err error) {
	g3.ZL().Error("connected",
		zap.String("uuid", conn.Uuid),
		zap.Reflect("query", conn.Query),
		zap.Error(err),
	)
	if conn.Conn != nil {
		_ = conn.Conn.Close()
	}
	if w.OnError != nil {
		w.OnError(conn, err)
	}
}

// HandleWebsocket HandleWebsocket
func HandleWebsocket(router string, opts ...WsWorkerOption) (*WsWorker, error) {
	if _, exist := workers[router]; exist {
		return nil, errors.New("Create Websocket Router Failed :" + router)
	}
	workers[router] = true
	worker := new(WsWorker)
	worker.Router = router
	worker.Upgrader = defaultUpgrader()
	worker.connections = make(map[string]*WsConn)
	for _, opt := range opts {
		opt(worker)
	}
	http.HandleFunc(router, func(writer http.ResponseWriter, request *http.Request) {
		var conn *websocket.Conn
		g3Conn := new(WsConn)
		g3Conn.Status = Connecting
		g3Conn.Uuid = fmt.Sprintf("%v", uuid.New())
		g3Conn.Query = helpers.ParseQueryString(request.RequestURI)
		g3Conn.CreateAt = time.Now()
		conn, err := worker.Upgrader.Upgrade(writer, request, nil)
		if err != nil {
			onError(worker, g3Conn, err)
			return
		}
		defer worker.closeConn(g3Conn)
		g3Conn.Conn = conn
		if !worker.handleConnect(g3Conn) {
			return
		}
		worker.listen(g3Conn)
	})
	return worker, nil
}
