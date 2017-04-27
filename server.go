package gowebsocket

import (
	"golang.org/x/net/websocket"

	"io"
	"net/http"
)

func NewWebsocket() WSServer {
	s := WSServer{}
	s.methods = make(map[string]MethodHandler)
	s.clients = make(map[string]WSClient)

	s.WS = &websocket.Server{Handler: s.handler, Handshake: s.handshake}

	return s
}

func (s *WSServer) Method(name string, fn MethodHandler) {
	s.methods[name] = fn
}

func (s *WSServer) handler(ws *websocket.Conn) {
	conn := WSConn{ws}
	defer ws.Close()

	for {
		msg, err := conn.ReadMessage()

		if err != nil {
			if err != io.EOF {
				// send error
				// Log.WithField("error", err).WithField("msg", msg).Error("Read Error")
			}
			break
		}

		s.handleMessage(&conn, &msg)
	}
}

func (s *WSServer) handleMessage(conn Connection, m *Message) {
	switch m.Msg {
	case "connect":
		s.handleConnect(conn, m)
	case "ping":
		s.handlePing(conn, m)
	case "method":
		s.handleMethod(conn, m)
	default:
		//send error
		// Log.WithField("msg", m).Error("Unknown Message Type")
		break
	}
}

func (s *WSServer) handleMethod(conn Connection, m *Message) {
	fn, ok := s.methods[m.Method]

	if !ok {
		// Log.WithField("method", m.Method).Error("Method Not Found")
		return
	}

	go fn(conn, m)
}

func (s *WSServer) handleConnect(conn Connection, m *Message) {
	id := "id" //RandomId(10)
	s.clients[id] = WSClient{id, conn}
	conn.WriteMessage(JsonData{"msg": "connected", "id": id})
}

func (s *WSServer) handleDisconnect() {
	delete(s.clients, "id")
}

func (s *WSServer) handlePing(conn Connection, m *Message) {
	msg := map[string]string{
		"msg": "pong",
	}

	if m.ID != "" {
		msg["id"] = m.ID
	}

	conn.WriteMessage(msg)
}

func (s *WSServer) handshake(config *websocket.Config, req *http.Request) error {
	// accept all connections
	return nil
}

func (s *WSServer) sendAll(msg *Message) {
	// for _, c := range s.clients {
	// c.Write(msg)
	// }
}
