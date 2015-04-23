package network

import (
	"database/sql"
	"errorHandler"
	"net"
	"settings"
)

type INetworkManager interface {
	Connect(ip string, remoteport int)
	Listen(localport int)
	Close(identifier ConnectionIdentifier)
	Write(identifier ConnectionIdentifier, msg string)
	Read(identifier ConnectionIdentifier)
	ConvertToConnectionIdentifier(ip string, localport, remoteport int) (identifier ConnectionIdentifier)
}

type ConnectionIdentifier struct {
	LocalAddress, RemoteAddress *net.TCPAddr
}

type NetworkManager struct {
	Database   *sql.DB
	Properties *settings.Settings
}

var tcpConnections map[ConnectionIdentifier]*net.TCPConn

func (networkManager *NetworkManager) Connect(ip string, remoteport int) {

}

func (networkManager *NetworkManager) Listen(localport int) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "")
	errorHandler.HandleError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	errorHandler.HandleError(err)
	tcpConnections[ConnectionIdentifier{conn.LocalAddr().(*net.TCPAddr), conn.RemoteAddr().(*net.TCPAddr)}] = conn
}

func (networkManager *NetworkManager) Close(identifier ConnectionIdentifier) {

}

func (networkManager *NetworkManager) Write(identifier ConnectionIdentifier, msg string) {

}

func (networkManager *NetworkManager) Read(identifier ConnectionIdentifier) {

}

func (networkManager *NetworkManager) ConvertToConnectionIdentifier(ip string, localport, remoteport int) (identifier ConnectionIdentifier) {

	return ConnectionIdentifier{}
}
