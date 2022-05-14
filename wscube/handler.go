package wscube

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const Version = "1"

type BusSubject string

type Handler struct {
	cubeInstance           Cube
	onlyAuthorizedRequests bool
	jwtSecret              string
	devMode                bool
	port                   int
	server                 *Server
	enableRouting          bool
	endpointsMap           map[Endpoint]Channel
	inputChannel           InputChannel
}

func parseEndpointsMap(rawMap string) (*map[Endpoint]Channel, error) {

	params := map[Endpoint]Channel{}

	if rawMap == "" {
		return &params, nil
	}

	for _, rawMap := range strings.Split(rawMap, ";") {
		splittedMap := strings.Split(rawMap, ":")

		if len(splittedMap) != 2 {
			return nil, fmt.Errorf("Wrong params format: %v\n", rawMap)
		}

		key := splittedMap[0]
		value := splittedMap[1]

		params[Endpoint(key)] = Channel(value)
	}

	return &params, nil
}

func (h *Handler) OnInitInstance() []InputChannel {
	return []InputChannel{
		InputChannel("wsinput"),
	}
}

func (h *Handler) OnStart(cubeInstance Cube) error {
	fmt.Println("Starting http gateway...")

	h.cubeInstance = cubeInstance
	h.jwtSecret = cubeInstance.GetParam("jwtSecret")
	h.onlyAuthorizedRequests = cubeInstance.GetParam("onlyAuthorizedRequests") == "true"
	h.devMode = cubeInstance.GetParam("dev") == "true"
	h.enableRouting = cubeInstance.GetParam("enableRouting") == "true"

	portString := cubeInstance.GetParam("port")

	var err error
	port := 80

	if portString != "" {
		port, err = strconv.Atoi(portString)
		if err != nil {
			cubeInstance.LogError("Wrong timeout")
			return err
		}
	}

	h.port = port

	endpointsMap, err := parseEndpointsMap(cubeInstance.GetParam("endpointsMap"))
	if err != nil {
		return err
	}

	h.endpointsMap = *endpointsMap

	h.server = NewServer(cubeInstance, h.devMode, h.enableRouting, *endpointsMap, h.onlyAuthorizedRequests, h.jwtSecret, port)
	go h.server.Start(cubeInstance)
	return nil
}

func (h *Handler) OnStop(c Cube) {
}

func (h *Handler) OnReceiveMessage(instance Cube, channel Channel, message Message) {

	switch message.Method {
	case "closeDeviceConnections":
		h.onCloseDeviceConnetions(message)
	case "closeUserConnections":
		h.onCloseUserConnetions(message)
	case "publishTextMessage":
		h.onSendMessage(message)

	default:
		fmt.Println("OnReceiveMessage: is not implemented")
		instance.LogError("OnReceiveMessage: is not implemented")
	}
}

func (h *Handler) onCloseDeviceConnetions(message Message) {

	if message.Params == nil {
		fmt.Println("onCloseDeviceConnetions: no params")
		return
	}

	var params CloseDeviceConnectionsParams
	err := json.Unmarshal(*message.Params, &params)
	if err == nil {
		fmt.Println("onCloseDeviceConnetions: wrong params")
		return
	}

	userId := (UserId)(params.UserId)
	deviceId := (DeviceId)(params.DeviceId)

	h.server.CloseDeviceConnections(userId, deviceId, params.Reason)
}

func (h *Handler) onCloseUserConnetions(message Message) {

	if message.Params == nil {
		fmt.Println("onCloseUserConnetions: no params")
		return
	}

	var params CloseUserConnectionsParams
	err := json.Unmarshal(*message.Params, &params)
	if err == nil {
		fmt.Println("onCloseUserConnetions: wrong params")
		return
	}

	userId := (UserId)(params.UserId)

	h.server.CloseUserConnections(userId, params.Reason)
}

func (h *Handler) onSendMessage(message Message) {

	if message.Params == nil {
		fmt.Println("onSendMessage: no params")
		return
	}

	var params PublishMessageParams
	err := json.Unmarshal(*message.Params, &params)
	if err == nil {
		fmt.Println("onSendMessage: wrong params")
		return
	}

	for _, receiver := range params.To {
		h.server.SendMessage(
			(*UserId)(receiver.UserId),
			(*DeviceId)(receiver.DeviceId),
			params.Type,
			params.Body,
		)
	}
}

// From bus
func (h *Handler) OnReceiveRequest(instance Cube, channel Channel, request Request) Response {
	fmt.Println("OnReceiveRequest: is not implemented")
	instance.LogError("OnReceiveRequest: is not implemented")
	return NewErrorResponse(
		"",
		"NotImplemented",
		"",
	)
}

var _ HandlerInterface = (*Handler)(nil)
