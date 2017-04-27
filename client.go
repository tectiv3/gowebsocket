package gowebsocket

import (
	"fmt"
	"io"
)

const channelBufSize = 100

var maxId int = 0

func NewClient(conn Connection, server *Server) *Client {
	if conn == nil {
		panic("conn cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, conn, server, ch, doneCh}
}

func (c *Client) Connection() Connection {
	return c.conn
}

func (c *Client) Send(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		c.server.Err(fmt.Errorf("client %d is disconnected.", c.id))
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

func (c *Client) Server() *Server {
	return c.server
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	for {
		select {
		// send message to the client
		case msg := <-c.ch:
			c.conn.WriteMessage(msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	for {
		select {
		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return
		// read data from websocket connection
		default:
			if msg, err := c.conn.ReadMessage(); err != nil {
				if err == io.EOF {
					c.doneCh <- true
				} else {
					c.server.Err(err)
				}
			} else {
				out := ClientMessage{c, &msg}
				c.server.Messages <- &out
			}
		}
	}
}
