package gowebsocket

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"sync"
)

var clientsMutex = sync.Mutex{}

func NewWebsocket(pattern string) *Server {
	return &Server{
		pattern:   pattern,
		clients:   make(map[int]*Client),
		addCh:     make(chan *Client),
		delCh:     make(chan *Client),
		sendAllCh: make(chan *Message),
		doneCh:    make(chan bool),
		errCh:     make(chan error),
		Messages:  make(chan *ClientMessage, channelBufSize),
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAll(msg *Message) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for _, client := range s.clients {
		client.Send(msg)
	}
}

func (s *Server) Clients() map[int]*Client {
	return s.clients
}

func (s *Server) Listen() {
	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			if err := ws.Close(); err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(&WebsocketConnection{ws}, s)
		s.Add(client)
		client.Listen()
	}

	http.Handle(s.pattern, websocket.Handler(onConnected))
	for {
		select {
		// Add new a client
		case c := <-s.addCh:
			clientsMutex.Lock()
			s.clients[c.id] = c
			log.Printf("Clients connected: %d", len(s.clients))
			clientsMutex.Unlock()

		// Del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			clientsMutex.Lock()
			delete(s.clients, c.id)
			clientsMutex.Unlock()

		// Broadcast message for all clients
		case msg := <-s.sendAllCh:
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
