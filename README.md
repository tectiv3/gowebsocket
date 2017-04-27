# gowebsocket

### Websocket server wrapper

#### Getting started

##### Requirements

- Go 1.7+ *earlier versions untested*

###### Instructions

1. Get deps:

```
go get github.com/tectiv3/gowebsocket
```

2. Example:

```go
package main

import (
	ws "github.com/tectiv3/gowebsocket"
	"log"
	"net/http"
	"os"
)

var (
	wsCh          ws.MessagesChannel
	wsServer      *ws.Server
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))

	wsServer = ws.NewWebsocket("/websocket")

	go func() {
		for {
			select {
			case incoming := <-wsServer.Messages:
				go answerMessage(incoming.Client, incoming.Msg)
			}
		}
	}()

	go wsServer.Listen()

	log.Println("Running webserver on port 8787")
	if err := http.ListenAndServe(":8787", nil); err != nil {
		Log.Error(err)
	}
}

func answerMessage(client *ws.Client, msg *ws.Message) {
	switch msg.Type {
	case "connect":
		client.Send(&ws.Message{Type: "success", Text: "connected"})
	}
}
```
