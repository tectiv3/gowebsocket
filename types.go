package gowebsocket

import (
	"golang.org/x/net/websocket"
)

const channelBufSize = 100

type Connection interface {
	ReadMessage() (Message, error)
	WriteMessage(*Message) error
}

type WebsocketConnection struct {
	ws *websocket.Conn
}

type Message struct {
	ID     int                    `json:"id,omitempty"`
	Text   string                 `json:"text,omitempty"`
	Type   string                 `json:"type"`
	Result interface{}            `json:"result,omitempty"`
	Params map[string]interface{} `json:"params,omitempty"`
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
