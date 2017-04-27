package gowebsocket

import (
	"golang.org/x/net/websocket"
)

type Connection interface {
	ReadMessage() (Message, error)
	WriteMessage(*Message) error
}

type WebsocketConnection struct {
	ws *websocket.Conn
}

type Message struct {
	ID     int                    `json:"id"`
	Text   string                 `json:"text"`
	Type   string                 `json:"type"`
	Result interface{}            `json:"result"`
	Params map[string]interface{} `json:"params"`
}

type ClientMessage struct {
	Client *Client
	Msg    *Message
}

type MessagesChannel <-chan *ClientMessage

type Server struct {
	pattern   string
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
	Messages  chan *ClientMessage
}

type Client struct {
	id     int
	conn   Connection
	server *Server
	ch     chan *Message
	doneCh chan bool
}

type JsonData map[string]interface{}
