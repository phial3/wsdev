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
	go s.handleMessages(conn)

}

func (s *Server) handleMessages(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		switch messageType {
		case websocket.TextMessage:
		case websocket.BinaryMessage:
		case websocket.CloseMessage:
			s.onMessage(messageType, conn, message)
			return
		}
	}
}

func (s *Server) onMessage(msgType int, conn *websocket.Conn, message []byte) {
	logger.Infof("message %s", message)
	conn.WriteMessage(msgType, message)
}
