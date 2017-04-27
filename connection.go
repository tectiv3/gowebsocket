package gowebsocket

import "golang.org/x/net/websocket"

func (c *WebsocketConnection) ReadMessage() (Message, error) {
	msg := Message{}
	err := websocket.JSON.Receive(c.ws, &msg)
	return msg, err
}

func (c *WebsocketConnection) WriteMessage(msg *Message) error {
	return websocket.JSON.Send(c.ws, msg)
}
