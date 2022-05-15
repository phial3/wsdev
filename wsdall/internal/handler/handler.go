package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wsdall/internal/logger"
	"github.com/wsdall/internal/tools"
	"github.com/wsdall/internal/wshub"
)

type Handler struct {
	wshub *wshub.WSHub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New() *Handler {
	return &Handler{
		wshub: wshub.New(),
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("connect error, ", err)
		return nil, errors.New("connect error")
	}

	// TODO: authenticate
	// _, err = h.Auth(r)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	h.wshub.Register("userId", ws)
	return ws, nil
}

// Auth is a middleware for websocket connection
func (h *Handler) Auth(req *http.Request) (string, error) {
	token := req.Header.Get("token")
	if token == "" {
		return "", errors.New("token is empty")
	}
	claims, err := tools.ParseToken(token)
	if err != nil {
		return "", errors.New("token is invalid")
	}
	return claims.Phone, nil
}
