package main

//
// import (
// 	"fmt"
// 	"log"
// 	"os"
//
// 	"github.com/urfave/cli"
// )
//
// func main() {
// 	app := cli.NewApp()
// 	app.Version = "0.0.3"
// 	app.Action = runServer
// 	app.Flags = []cli.Flag{
// 		cli.StringFlag{
// 			Name:   "bus-host",
// 			EnvVar: "GATEWAY_BUS_HOST",
// 			Usage:  "bus host",
// 		},
// 		cli.IntFlag{
// 			Name:   "bus-port",
// 			EnvVar: "GATEWAY_BUS_PORT",
// 			Usage:  "bus port",
// 		},
// 		cli.StringFlag{
// 			Name:   "jwt-secret",
// 			EnvVar: "GATEWAY_JWT_SECRET",
// 			Usage:  "jwt secret",
// 		},
// 		cli.IntFlag{
// 			Name:   "max-connections",
// 			EnvVar: "GATEWAY_MAX_CONNECTIONS",
// 			Usage:  "maximum number of connections",
// 		},
// 		cli.StringFlag{
// 			Name:   "endpoints-map",
// 			EnvVar: "GATEWAY_ENDPOINTS_MAP",
// 			Usage:  "map endpoint to channel",
// 		},
// 		cli.StringFlag{
// 			Name:   "input-channel",
// 			EnvVar: "GATEWAY_INPUT_CHANNEL",
// 			Usage:  "input channel, default \"wsinput\"",
// 		},
// 		cli.BoolTFlag{
// 			Name:   "only-authorized-requests",
// 			EnvVar: "GATEWAY_ONLY_AUTHORIZED_REQUESTS",
// 			Usage:  "handle only authorized requests",
// 		},
// 		cli.BoolFlag{
// 			Name:   "enable-routing",
// 			EnvVar: "GATEWAY_ENABLE_ROUTING",
// 		},
// 		cli.BoolFlag{
// 			Name:   "dev",
// 			EnvVar: "GATEWAY_DEV",
// 			Usage:  "log all requests",
// 		},
// 		cli.StringFlag{
// 			Name:   "port",
// 			EnvVar: "GATEWAY_PORT",
// 			Usage:  "port to listen",
// 		},
// 	}
//
// 	err := app.Run(os.Args)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// func runServer(c *cli.Context) error {
//
// 	busHost := c.String("bus-host")
// 	if busHost == "" {
// 		return fmt.Errorf("bus host is required")
// 	}
//
// 	busPort := c.Int("bus-port")
// 	if busPort == 0 {
// 		return fmt.Errorf("bus port is required")
// 	}
//
// 	jwtSecret := c.String("jwt-secret")
// 	if jwtSecret == "" {
// 		return fmt.Errorf("jwt secret is required")
// 	}
//
// 	maxConnections := c.String("max-connections")
// 	if maxConnections == "" {
// 		return fmt.Errorf("max connections is required")
// 	}
//
// 	_ = c.String("port")
//
// 	// onlyAuthorizedRequests := "true"
// 	// if c.Bool("only-authorized-requests") {
// 	// 	onlyAuthorizedRequests = "true"
// 	// } else {
// 	// 	onlyAuthorizedRequests = "false"
// 	// }
// 	//
// 	// dev := "false"
// 	// if c.Bool("dev") {
// 	// 	dev = "true"
// 	// } else {
// 	// 	dev = "false"
// 	// }
//
// 	// enableRouting := "false"
// 	// endpointsMap := c.String("endpoints-map")
// 	//
// 	// if c.Bool("enable-routing") {
// 	// 	enableRouting = "true"
// 	//
// 	// 	if endpointsMap == "" {
// 	// 		return fmt.Errorf("endpoints map is required")
// 	// 	}
// 	//
// 	// } else {
// 	// 	enableRouting = "false"
// 	// }
//
// 	// channelsMapping := map[cube_executor.CubeChannel]cube_executor.BusChannel{}
// 	// inputChannel := c.String("input-channel")
// 	// if inputChannel != "" {
// 	// 	channelsMapping[cube_executor.CubeChannel("wsinput")] = cube_executor.BusChannel(inputChannel)
// 	// }
// 	//
// 	// cube, err := cube_executor.NewCube(cube_executor.CubeConfig{
// 	// 	BusPort:         busPort,
// 	// 	BusHost:         busHost,
// 	// 	ChannelsMapping: channelsMapping,
// 	// 	Params: map[string]string{
// 	// 		"jwtSecret":              jwtSecret,
// 	// 		"maxConnections":         maxConnections,
// 	// 		"endpointsMap":           endpointsMap,
// 	// 		"onlyAuthorizedRequests": onlyAuthorizedRequests,
// 	// 		"dev":                    dev,
// 	// 		"port":                   port,
// 	// 		"enableRouting":          enableRouting,
// 	// 	},
// 	// }, &cube_websocket_gateway.Handler{})
// 	//
// 	// if err != nil {
// 	// 	return fmt.Errorf("can't start: %v", err)
// 	// }
//
// 	// return cube.Start()
// 	return nil
// }
