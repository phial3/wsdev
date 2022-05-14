package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wsdall/internal/wshub"
)

type Handler struct {
	wshub *wshub.WSHub
}

func New() *Handler {
	return &Handler{
		wshub: wshub.New(),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// TODO: add autentication
	usn := r.Header.Get("username")
	h.wshub.Register(ws, usn)
}
