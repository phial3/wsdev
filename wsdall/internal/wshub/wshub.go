package wshub

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	logger "github.com/wsdall/internal"
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

func (w *WSHub) Register(conn *websocket.Conn, key string) {
	w.connMap.Store(key, conn)
	go w.listen(conn)
}

func (w *WSHub) Unregister(key string) {
	conn, loaded := w.connMap.LoadAndDelete(key)
	if loaded {
		conn.(*websocket.Conn).Close()
	}
}

// proccess the message
// locking
func (w *WSHub) Routes() {
	for {
		select {
		case msgB := <-w.rChan:
			// logrus.Debugf("outgoing message: %s", string(msgB))
			logger.GetLogger().Debugf("outgoing message: %s", string(msgB))
			msg := models.Message{}
			json.Unmarshal(msgB, &msg)
			w.write(msg.Receiver, msgB)
		}
	}
}

func (w *WSHub) write(key string, msg []byte) {
	c, loaded := w.connMap.Load(key)
	if !loaded {
		// logrus.Warn("unable to find connection")
		logger.GetLogger().Warn("unable to find connection")
		return
	}
	conn := c.(*websocket.Conn)
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
		// logrus.Debug("stop listening")
		logger.GetLogger().Debug("stop listening")
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// logrus.Debugf("incomming message: %s", string(p))
		logger.GetLogger().Debugf("receive message: %s", string(p))
		w.rChan <- p
	}
}
