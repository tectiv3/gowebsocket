package gowebsocket

import (
	"golang.org/x/net/websocket"
)

type Connection interface {
	ReadMessage() (Message, error)
	WriteMessage(interface{}) error
}

type WSConn struct {
	ws *websocket.Conn
}

type Message struct {
	ID     string                 `json:"id"`
	Msg    string                 `json:"msg"`
	Method string                 `json:"method"`
	Result interface{}            `json:"result"`
	Params map[string]interface{} `json:"params"`
}

type MethodHandler func(Connection, *Message)

type WSServer struct {
	methods map[string]MethodHandler
	clients map[string]WSClient
	WS      *websocket.Server
}

type WSClient struct {
	Id  string
	Con Connection
}

type JsonData map[string]interface{}
