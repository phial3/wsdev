package bitmex

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wsbitmex/pkg/config"
	t "github.com/wsbitmex/pkg/types"
)

type StockMarket struct {
	conn *websocket.Conn
	Msgs chan []byte
	conf *config.StockMarket
}

type Bitmex interface {
	AuthKeyExpires() (err error)
	Subscribe() (err error)
	Run(verbose bool)
}

var _ Bitmex = (*StockMarket)(nil)

func (b *StockMarket) AuthKeyExpires() (err error) {

	_, message, err := b.conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("->:: %s", message)
	auth_json := fmt.Sprintf(`{"op": "authKeyExpires", "args": ["%s", %d , "%s"]}`,
		b.conf.ApiKey,
		b.conf.Expired,
		b.conf.Signature)
	err = b.conn.WriteMessage(websocket.TextMessage, []byte(auth_json))
	if err != nil {
		log.Fatal("auth:", err)
	}
	_, message, err = b.conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("->:: %s", message)
	resp := t.ResponseAuth{}
	err = json.Unmarshal([]byte(message), &resp)
	if err != nil {
		log.Println("unmarshal:", err)
		return
	}
	if !resp.Success {
		return errors.New("auth: not success")
	}
	return

}

func (b *StockMarket) Subscribe() (err error) {

	sub_json := fmt.Sprintf(`{"op": "subscribe", "args": ["%s"]}`,
		b.conf.Topic)
	err = b.conn.WriteMessage(websocket.TextMessage, []byte(sub_json))
	if err != nil {
		log.Fatal("subscribe:", err)
		return
	}
	return
}

func (b *StockMarket) Run(verbose bool) {
	go func() {
		for {
			_, message, err := b.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			if verbose {
				log.Printf("->: %s", message)
			}
			resp := &t.Response{}
			err = json.Unmarshal(message, resp)
			if err != nil {
				log.Println("unmarshal:", err)
				continue
			}
			for _, data := range resp.Data {

				// TODO: check empty field
				if data.Price == 0 {
					err := errors.New("data: no lastPrice")
					if verbose {
						log.Println(err)
					}
					continue
				}
				clientResponse := &t.ClientResponse{
					Timestamp: data.Timestamp,
					Symbol:    data.Symbol,
					Price:     data.Price,
				}
				respMessage, err := json.Marshal(clientResponse)
				if err != nil {
					log.Println("marshal:", err)
					continue
				}
				if verbose {
					log.Printf("<-: %s", respMessage)
				}
				b.Msgs <- respMessage
			}
		}
	}()
}

func websocketApiKeyAuth(apiSecret string, path string, expired int) (signature string) {

	// TODO

	log.Println("websocketApiKeyAuth: unrealised; default signature is used")
	return "b4f9ff9549ba051c3558a8caf8e6f41bd9fbeba0fd734e9919e0da7353574cd6"
}

func NewBitmexConfig(address string,
	apiKey string,
	apiSecret string,
	expired int) *config.StockMarket {
	return &config.StockMarket{
		Scheme:    "wss",
		Address:   address,
		Path:      "/realtime",
		ApiKey:    apiKey,
		Expired:   expired,
		Signature: websocketApiKeyAuth(apiSecret, "/realtime", expired),
		Topic:     "instrument",
	}
}

func NewStockMarket(conn *websocket.Conn,
	msgs chan []byte,
	conf *config.StockMarket) *StockMarket {
	return &StockMarket{
		conn: conn,
		Msgs: msgs,
		conf: conf,
	}
}
