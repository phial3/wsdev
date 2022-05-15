package wscube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type Endpoint string

type Server struct {
	cubeInstance           Cube
	upgrade                websocket.Upgrader
	devMode                bool
	httpServer             *http.Server
	onlyAuthorizedRequests bool
	jwtSecret              string
	connections            *ConnectionManager
	lastConnectionNumber   int64
	port                   int
	enableRouting          bool
	endpointsMap           map[Endpoint]Channel
}

func NewServer(
	cubeInstance Cube,
	devMode bool,
	enableRouting bool,
	endpointsMap map[Endpoint]Channel,
	onlyAuthorizedRequests bool,
	jwtSecret string, port int) *Server {
	return &Server{
		cubeInstance:           cubeInstance,
		upgrade:                websocket.Upgrader{},
		devMode:                devMode,
		onlyAuthorizedRequests: onlyAuthorizedRequests,
		jwtSecret:              jwtSecret,
		connections:            NewConnectionsStorage(),
		port:                   port,
		enableRouting:          enableRouting,
		endpointsMap:           endpointsMap,
	}
}

func (s *Server) Start(cubeInstance Cube) {

	srv := http.Server{
		Addr:    ":" + strconv.Itoa(s.port),
		Handler: s,
	}

	s.httpServer = &srv

	fmt.Println("Start http listening")
	cubeInstance.LogInfo("Start http listening")

	address := fmt.Sprintf(":%v", s.port)

	http.HandleFunc("/", s.ServeHTTP)
	err := http.ListenAndServe(address, nil)

	fmt.Println("Stop http listenning", err)
	cubeInstance.LogFatal(err.Error())
}

func (s *Server) getAuthData(tokenString string) (*UserId, *DeviceId, error) {

	if tokenString == "" {
		return nil, nil, fmt.Errorf("empty token")
	}

	// newToken, err := jws.ParseJWT([]byte(tokenString))
	// if err != nil {
	// 	return nil, nil, err
	// }
	//
	// err = newToken.Validate([]byte(s.jwtSecret), crypto.SigningMethodHS512)
	// if err != nil {
	// 	return nil, nil, err
	// }
	//
	// claims := newToken.Claims()
	// userId := UserId(claims.Get("userId").(string))
	// deviceId := DeviceId(claims.Get("deviceId").(string))

	return nil, nil, nil
}

// On connection
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var userId *UserId
	var deviceId *DeviceId
	var err error

	if s.devMode {
		fmt.Println("")
		fmt.Println("-----")
		fmt.Println("RECEIVE REQUEST:")
		fmt.Println("method: ", request.Method)
		fmt.Println("url: ", request.URL)
		fmt.Println("uri: ", request.RequestURI)
		fmt.Println("headers: ", request.Header)
		fmt.Println("body:")
		fmt.Println(request.Body)
		fmt.Println("-----")
	}

	secData := request.Header["Sec-Websocket-Protocol"]
	if secData == nil || len(secData) != 1 || !strings.HasPrefix(secData[0], "token,") {
		http.Error(writer,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(secData[0], "token,")
	token = strings.TrimSpace(token)

	if token != "" && s.jwtSecret != "" {

		userId, deviceId, err = s.getAuthData(token)

		if err != nil {
			http.Error(writer,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
	}

	if s.onlyAuthorizedRequests && userId == nil {
		http.Error(writer,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
		return
	}

	responseHeader := http.Header{}
	responseHeader.Set("Sec-WebSocket-Protocol", "token")
	connection, err := s.upgrade.Upgrade(writer, request, responseHeader)
	if err != nil {
		return
	}

	fmt.Println(deviceId)
	connection.SetReadLimit(100000000)
	con := s.registerConnection(connection)
	// TODO: add onlyAuthorized connections support
	con.Login(*userId, *deviceId)

	go s.handleInputMessages(con)
	s.cleanConnectionsIfNeed()

	packedMessage, _ := s.packMessage(userId, deviceId, "onConnect", &[]byte{})
	s.cubeInstance.PublishMessage(Channel("wsOutput"), *packedMessage)

	// TODO: add onlyAuthorized connections support
}

func (s *Server) cleanConnectionsIfNeed() {

	now := time.Now().Unix()
	cnt := s.connections.GetCountConnection()
	if cnt > 200 {
		s.connections.RemoveIf(func(con *Connection) bool {

			return now-con.GetStartTime().Unix() > 60

		}, func(connections []*Connection) {

			for _, connection := range connections {
				connection.Close(websocket.ClosePolicyViolation, "Auth")
			}
		})
	}
}

func (s *Server) handleInputMessages(netConnection *Connection) {

	for {
		messageType, message, err := netConnection.ReadMessage()
		if err != nil {
			netConnection.Close(websocket.CloseInternalServerErr, "ServerError")
			s.onClose(netConnection)
			return
		}

		netConnection.UpdateLastPingTime()

		switch messageType {
		case websocket.TextMessage:
			s.onReceiveMessage(netConnection, true, &message)
		case websocket.BinaryMessage:
			s.onReceiveMessage(netConnection, false, &message)
		case websocket.CloseMessage:
			s.onClose(netConnection)
			return
		}
	}
}

func (s *Server) getNewConnectionId() ConnectionId {
	return ConnectionId(atomic.AddInt64(&s.lastConnectionNumber, 1))
}

func (s *Server) registerConnection(connection *websocket.Conn) *Connection {

	wsConnection := NewConnection(s.getNewConnectionId(), connection)
	s.connections.AddNewConnection(wsConnection)

	connection.SetCloseHandler(func(code int, text string) error {
		s.onClose(wsConnection)
		return nil
	})

	return wsConnection
}

func (s *Server) unregisterConnection(connection *Connection) {
	s.connections.RemoveConnection(connection)
}

func (s *Server) onClose(connection *Connection) {

	connectionId, userId, deviceId := connection.GetInfo()
	if connectionId == -1 {
		return
	}

	s.unregisterConnection(connection)
	packedMessage, _ := s.packMessage(&userId, &deviceId, "onClose", &[]byte{})
	s.cubeInstance.PublishMessage(Channel("wsOutput"), *packedMessage)
}

func (s *Server) onReceiveMessage(connection *Connection, isText bool, rawBody *[]byte) {

	outputChannel := Channel("wsOutput")
	body := rawBody

	if s.enableRouting {

		var packet RoutingPacket
		err := json.Unmarshal(*rawBody, &packet)
		if err != nil {
			connection.SendText([]byte("ErrorParsingRoutingPacket"))
			return
		}

		outputChannel = s.endpointsMap[Endpoint(packet.Endpoint)]
		if outputChannel == "" {
			connection.SendText([]byte("ErrorEndpointNotFound"))
			return
		}

		if len(packet.Payload) == 0 {
			connection.SendText([]byte("ErrorEmptyPayload"))
			return
		}

		body = (*[]byte)(&packet.Payload)

	} else {
		mapChannel := s.endpointsMap["wsOutput"]
		if mapChannel != "" {
			outputChannel = mapChannel
		}
	}

	method := "onTextMessage"
	if !isText {
		method = "onBinaryMessage"
	}

	_, userId, deviceId := connection.GetInfo()
	packedMessage, err := s.packMessage(&userId, &deviceId, method, body)
	if err != nil {
		return
	}

	s.cubeInstance.PublishMessage(outputChannel, *packedMessage)
}

func (s *Server) packMessage(userId *UserId, deviceId *DeviceId, method string, body *[]byte) (*Message, error) {

	params := OnReceiveMessageParams{
		DeviceId:  (*string)(deviceId),
		UserId:    (*string)(userId),
		InputTime: time.Now().UnixNano(),
		Body:      *body,
	}

	packedParams, _ := json.Marshal(params)

	messageData := &Message{
		Method: method,
		Params: (*json.RawMessage)(&packedParams),
	}

	return messageData, nil
}

func (s *Server) CloseDeviceConnections(userId UserId, deviceId DeviceId, reason string) {
	s.connections.RemoveDeviceConnections(userId, deviceId, func(connections []*Connection) {

		for _, connection := range connections {
			connection.Close(websocket.CloseNormalClosure, reason)
		}
	})
}

func (s *Server) CloseUserConnections(userId UserId, reason string) {
	s.connections.RemoveUserConnections(userId, func(connections []*Connection) {

		for _, connection := range connections {
			connection.Close(websocket.CloseNormalClosure, reason)
		}
	})
}

func (s *Server) SendMessage(userId *UserId, deviceId *DeviceId, messageType MessageType, message []byte) {
	connections := []*Connection{}
	if deviceId != nil {
		connections = s.connections.GetDeviceConnections(*userId, *deviceId)
	} else if userId != nil {
		connections = s.connections.GetUserConnections(*userId)
	}

	if len(connections) > 0 {
		for _, connection := range connections {
			switch messageType {
			case TEXT:
				connection.SendText(message)
			case BINARY:
				connection.SendBinary(message)
			}

		}
	}
}
