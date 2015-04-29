package network

import (
	"database/sql"
	"errorHandler"
	"net"
	"settings"
	"strconv"
)

type INetworkManager interface {
	Connect(ip string, remoteport int)
	Listen(localport int)
	Close(identifier string)
	Write(identifier string, msg string)
	Read(identifier string)
	ConvertToIdentifier(ip string, localport, remoteport int) (identifier string)
	ConvertToStrings(identifier string) (ip, localport, remoteport string)
}

type NetworkManager struct {
	Database   *sql.DB
	Properties *settings.Settings
}

var tcpConnections map[string]*net.TCPConn

func (networkManager *NetworkManager) Connect(ip string, remoteport int) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+strconv.Itoa(remoteport))
	errorHandler.HandleError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	errorHandler.HandleError(err)
	tcpConnections[convertToIdentifier(conn.LocalAddr().(*net.TCPAddr), conn.RemoteAddr().(*net.TCPAddr))] = conn
}

func (networkManager *NetworkManager) Listen(localport int) {
}

func (networkManager *NetworkManager) Close(identifier string) {

}

func (networkManager *NetworkManager) Write(identifier string, msg string) {
	/*conn, contained := tcpConnections[identifier]
	
	if contained {
		//conn.Write(msg)
	}*/
}

func (networkManager *NetworkManager) Read(identifier string) {

}

func convertToIdentifier(localAddr, remoteAddr *net.TCPAddr) (identifier string) {

	return ""
}

func (networkManager *NetworkManager) ConvertToIdentifier(ip string, localport, remoteport int) (identifier string) {

	return ""
}

func (networkManager *NetworkManager) ConvertToStrings(identifier string) (ip, localport, remoteport string) {

	return "", "", ""
}
