package main

import (
	"flag"
	"fmt"
	"net/http"
)

import (
	"golang.org/x/net/context"
)

var (
	grpcAddr  = flag.String("grpcaddr", ":8001", "listen grpc addr")
	httpAddr  = flag.String("addr", ":8000", "listen http addr")
	debugAddr = flag.String("debugaddr", ":8002", "listen debug addr")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go http.ListenAndServe(*debugAddr, nil)
	fmt.Println("listening")
	// http.ListenAndServe(*httpAddr, proxy.WebsocketProxy())
	return nil
}
