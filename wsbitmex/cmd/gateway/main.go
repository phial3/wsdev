package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wsbitmex/pkg/bitmex"
	bit "github.com/wsbitmex/pkg/bitmex"
	"github.com/wsbitmex/pkg/config"
	sub "github.com/wsbitmex/pkg/subscriber"
	// TODO: add kitlog "github.com/go-kit/kit/log"
)

var addr = flag.String("addr", "testnet.bitmex.com", "http service address")
var apiKey = flag.String("apikey", "ORqVaoVf1TJrVnKexpWjHfjk", "")
var apiSecret = flag.String("apisecret", "", "")
var expired = flag.Int("expired", 1640984399, "timestamp")
var debug = flag.Bool("debug", false, "")

func createConn(conf *config.StockMarket) (c *websocket.Conn, err error) {
	u := url.URL{
		Scheme: conf.Scheme,
		Host:   conf.Address,
		Path:   conf.Path,
	}
	log.Printf("connecting to %s", u.String())

	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	return
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	smConf := bit.NewBitmexConfig(*addr, *apiKey, *apiSecret, *expired)

	c, _ := createConn(smConf)
	defer c.Close()

	ch := make(chan []byte, 100)
	bitmex := bitmex.NewStockMarket(c, ch, smConf)

	err := bitmex.AuthKeyExpires()
	if err != nil {
		log.Fatal("auth:", err)
	}

	err = bitmex.Subscribe()
	if err != nil {
		log.Fatal("sub:", err)
	}

	gwConf := sub.NewGatewayConfig("0.0.0.0", 9090)
	bindAddr := fmt.Sprintf("%s:%d", gwConf.Host, gwConf.Port)

	cl := make(map[string]*websocket.Conn)
	gateway := &sub.Gateway{Clients: cl}
	http.HandleFunc("/", gateway.Subscribe)

	go func() {
		log.Fatal(http.ListenAndServe(bindAddr, nil))
	}()

	bitmex.Run(*debug)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-bitmex.Msgs:
			msg, ok := <-bitmex.Msgs
			if !ok {
				break
			}
			_ = gateway.Publish(msg)
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}

}
