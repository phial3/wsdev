package main

import (
	logger "github.com/wsdall/internal"
	"github.com/wsdall/internal/server"
)

func main() {
	logger.GetLogger().Info("Starting ws gateway server...")
	srv := server.New()
	srv.Start()
}
