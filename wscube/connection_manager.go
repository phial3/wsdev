package wscube

import (
	"sync"
)

type ConnectionManager struct {
	mutex                        sync.RWMutex
	connectionsById              map[ConnectionId]*Connection
	numberOfNotLoggedConnections int
}

func NewConnectionsStorage() *ConnectionManager {
	return &ConnectionManager{
		mutex:                        sync.RWMutex{},
		connectionsById:              make(map[ConnectionId]*Connection),
		numberOfNotLoggedConnections: 0,
	}
}

func (s *ConnectionManager) AddNewConnection(connection *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.numberOfNotLoggedConnections++
	s.connectionsById[connection.id] = connection
}

func (s *ConnectionManager) RemoveConnection(connection *Connection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.removeConnection(connection)
}

func (s *ConnectionManager) removeConnection(connection *Connection) {

	connectionId, userId, _ := connection.GetInfo()

	connectionBefore := s.connectionsById[connectionId]
	if connectionBefore == nil {
		return
	}

	delete(s.connectionsById, connectionId)

	if userId == "" {
		s.numberOfNotLoggedConnections--
		return
	}
}

func (s *ConnectionManager) GetUserConnections(userId UserId) []*Connection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	connections := []*Connection{}

	for _, connection := range s.connectionsById {
		if connection.userId == userId {
			connections = append(connections, connection)
		}
	}

	return connections
}

func (s *ConnectionManager) GetDeviceConnections(userId UserId, deviceId DeviceId) []*Connection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	connections := []*Connection{}

	for _, connection := range s.connectionsById {
		if connection.deviceId == deviceId && connection.userId == userId {
			connections = append(connections, connection)
		}
	}

	return connections
}

func (s *ConnectionManager) GetConnectionById(connectionId ConnectionId) *Connection {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.connectionsById[connectionId]
}

func (s *ConnectionManager) GetCountConnection() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.numberOfNotLoggedConnections
}

func (s *ConnectionManager) RemoveIf(condition func(con *Connection) bool, afterRemove func(connections []*Connection)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	connections := []*Connection{}

	for id, connection := range s.connectionsById {
		if condition(connection) {
			delete(s.connectionsById, id)
			connections = append(connections, connection)
		}
	}

	afterRemove(connections)
}

func (s *ConnectionManager) RemoveDeviceConnections(userId UserId, deviceId DeviceId, afterRemove func(connections []*Connection)) {
	s.RemoveIf(func(con *Connection) bool {
		return con.deviceId == deviceId && con.userId == userId
	}, afterRemove)
}

func (s *ConnectionManager) RemoveUserConnections(userId UserId, afterRemove func(connections []*Connection)) {
	s.RemoveIf(func(con *Connection) bool {
		return con.userId == userId
	}, afterRemove)
}
