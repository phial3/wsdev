package main

import (
	"github.com/wsdall/internal/config"
	"github.com/wsdall/internal/logger"
	"github.com/wsdall/internal/server"
)

func main() {
	logger.Info("Starting ws gateway server...")
	srv := server.New()
	config.LoadConfig()
	srv.Start()
}
