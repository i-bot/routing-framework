package network

import (
	"database/sql"
	"net"
	"errorHandler"
	"settings"
)

type INetworkManager interface {
	ConnectTo(addr string)
}

type ConnectionIdentifier struct {
	LocalAddress, RemoteAddress *net.TCPAddr
}

type NetworkManager struct{
	Database *sql.DB
	Properties *settings.Settings
}

var tcpConnections map[ConnectionIdentifier]*net.TCPConn

func (networkManager *NetworkManager) ConnectTo(addr string){
}
