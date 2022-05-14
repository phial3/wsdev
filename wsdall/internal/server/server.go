package server

import (
	"log"
	"net/http"

	"github.com/wsdall/internal/handler"
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
	http.HandleFunc("/wsconnect", s.handler.Connect)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
