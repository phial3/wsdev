package subscriber

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/wsbitmex/pkg/config"
	t "github.com/wsbitmex/pkg/types"
)

type Subscriber interface {
	Subscribe(w http.ResponseWriter, r *http.Request)
	Publish(msg []byte) (err error)
}

var _ Subscriber = (*Gateway)(nil)

type Gateway struct {
	Clients  map[string]*websocket.Conn
	mux      sync.Mutex
	upgrader websocket.Upgrader
}

// TODO: add subscribe to selected topics
func (g *Gateway) Subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var resp string
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		connId := r.Header.Get("Sec-Websocket-Key")
		req := &t.ClientRequest{}
		err = json.Unmarshal(message, req)
		if err != nil {
			log.Println("unmarshal:", err)
			continue
		}
		if req.Action == "subscribe" {
			resp = "subscribed"
			g.add(connId, c)

		}
		if req.Action == "unsubscribe" {
			resp = "unsubscribed"
			g.delete(connId)
		}

		log.Printf("recv: %s, %s", connId, string(resp))
		err = c.WriteMessage(websocket.TextMessage, []byte(resp))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (g *Gateway) Publish(msg []byte) (err error) {
	g.mux.Lock()
	defer g.mux.Unlock()
	var wg sync.WaitGroup
	for id, c := range g.Clients {
		wg.Add(1)
		id, c := id, c
		go func() {
			defer wg.Done()
			_ = writeMessage(msg, id, c)
		}()

	}
	wg.Wait()
	return
}

func (g *Gateway) add(id string, c *websocket.Conn) {
	g.mux.Lock()
	g.Clients[id] = c
	g.mux.Unlock()
}

func (g *Gateway) delete(id string) {
	g.mux.Lock()
	delete(g.Clients, id)
	g.mux.Unlock()
}

func writeMessage(msg []byte, id string, c *websocket.Conn) (err error) {
	err = c.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("write:", err)
		return
	}
	log.Printf("recv to %s: %s", id, string(msg))
	return
}

func NewGatewayConfig(host string,
	port int) *config.Gateway {
	return &config.Gateway{
		Host: host,
		Port: port,
	}
}
