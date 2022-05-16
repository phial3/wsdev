package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wsdall/internal/handler"
	"github.com/wsdall/internal/logger"
)

type Server struct {
	handler *handler.Handler
}

func New() *Server {
	return &Server{
		handler: handler.New(),
	}
}

func (s *Server) Start() {
	port := 8080
	address := fmt.Sprintf(":%v", port)
	log.Printf("Websocket Server started at %s", address)
	http.HandleFunc("/", s.ServeHTTP)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("Stop Websocket Server", err)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := s.handler.Connect(w, req)
	if err == nil {
		logger.Errorf("connection error %v", err)
		return
	}
	go s.onMessage(conn)
}

func (s *Server) onConnected(conn *websocket.Conn) {

}

func (s *Server) onMessage(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("read msg error %v", err)
			conn.Close()
			return
		}

		switch messageType {
		case websocket.TextMessage:
			conn.WriteMessage(websocket.TextMessage, message)
		case websocket.BinaryMessage:
			conn.WriteMessage(websocket.BinaryMessage, message)
		case websocket.CloseMessage:
			conn.Close()
		default:
			logger.Warnf("unknown message type %v", messageType)
			return
		}
	}
}

func (s *Server) onClose(conn *websocket.Conn) {

}
