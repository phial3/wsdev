package wshub

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/wsdall/internal/logger"
	"github.com/wsdall/internal/models"
)

type WSHub struct {
	connMap *sync.Map
	rChan   chan []byte
}

func New() *WSHub {
	w := &WSHub{
		connMap: &sync.Map{},
		rChan:   make(chan []byte),
	}
	go w.Routes()
	return w
}

func (w *WSHub) Register(key string, conn *websocket.Conn) {
	w.connMap.Store(key, conn)
	go w.listen(conn)
}

func (w *WSHub) Unregister(key string) {
	conn, ok := w.connMap.LoadAndDelete(key)
	if ok {
		conn.(*websocket.Conn).Close()
	}
}

// proccess the message
// locking
func (w *WSHub) Routes() {
	for {
		select {
		case msgB := <-w.rChan:
			logger.Debugf("outgoing message: %s", string(msgB))
			msg := models.Message{
				Sender:   "",
				Receiver: "userId",
				Content:  "",
			}
			json.Unmarshal(msgB, &msg)
			w.write(msg.Receiver, msgB)
		}
	}
}

func (w *WSHub) write(key string, msg []byte) {
	data, ok := w.connMap.Load(key)
	if !ok {
		logger.Warn("unable to find connection")
		return
	}
	conn := data.(*websocket.Conn)
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println(err)
		return
	}
}

// listen on websocket message
// closed when the connection is closed
// locking
func (w *WSHub) listen(conn *websocket.Conn) {
	defer func() {
		logger.Debug("stop listening")
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		logger.Debugf("receive message: %s", string(data))
		w.rChan <- data
	}
}
