package network

import (
	"bytes"
	"database/sql"
	"errorHandler"
	"fmt"
	"io"
	"net"
	"settings"
	"strconv"
	"strings"
)

type INetworkManager interface {
	Init()
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

func (networkManager *NetworkManager) Init() {
	tcpConnections = make(map[string]*net.TCPConn)
}

func (networkManager *NetworkManager) Connect(ip string, remoteport int) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+strconv.Itoa(remoteport))
	errorHandler.HandleError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err == nil {
		identifier := convertToIdentifier(conn.LocalAddr().(*net.TCPAddr), conn.RemoteAddr().(*net.TCPAddr))

		tcpConnections[identifier] = conn
		HandleOpen(networkManager, identifier)

		networkManager.Read(identifier)
	} else {
		fmt.Println("Connect(): " + err.Error())
	}
}

func (networkManager *NetworkManager) Listen(localport int) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(localport))
	errorHandler.HandleError(err)

	listener, err := net.ListenTCP("tcp4", tcpAddr)
	errorHandler.HandleError(err)

	accept := func() {
		for {
			conn, err := listener.AcceptTCP()
			if err == nil {
				identifier := convertToIdentifier(conn.LocalAddr().(*net.TCPAddr), conn.RemoteAddr().(*net.TCPAddr))

				tcpConnections[identifier] = conn
				HandleOpen(networkManager, identifier)

				networkManager.Read(identifier)
			} else {
				fmt.Println("Listen(): " + err.Error())
				break
			}
		}
	}

	go accept()
}

func (networkManager *NetworkManager) Close(identifier string) {
	conn, available := tcpConnections[identifier]
	if available {
		delete(tcpConnections, identifier)
		err := conn.Close()
		if err == nil {
			HandleClose(networkManager, identifier)
		} else {
			fmt.Println("Close(): " + err.Error())
		}
	}
}

func (networkManager *NetworkManager) Write(identifier string, msg string) {
	conn, available := tcpConnections[identifier]
	if available {
		buf := []byte(msg)
		_, err := conn.Write(buf)
		if err == nil {
			HandleWrite(msg, networkManager, identifier)
		} else {
			fmt.Println("Write(): " + err.Error())
		}
	}
}

func (networkManager *NetworkManager) Read(identifier string) {
	read := func() {
		conn, available := tcpConnections[identifier]
		if available {
			for {
				var buf bytes.Buffer
				_, err := io.Copy(&buf, conn)
				if err == nil {
					HandleRead(buf.String(), networkManager, identifier)
				} else {
					fmt.Println("Read(): " + err.Error())
					break
				}
			}
		}
	}

	go read()
}

func convertToIdentifier(localAddr, remoteAddr *net.TCPAddr) (identifier string) {
	return remoteAddr.IP.String() + ":" + strconv.Itoa(localAddr.Port) + ":" + strconv.Itoa(remoteAddr.Port)
}

func (networkManager *NetworkManager) ConvertToIdentifier(ip string, localport, remoteport int) (identifier string) {
	return ip + ":" + strconv.Itoa(localport) + ":" + strconv.Itoa(remoteport)
}

func (networkManager *NetworkManager) ConvertToStrings(identifier string) (ip, localport, remoteport string) {
	split := strings.Split(identifier, ":")
	return split[0], split[1], split[2]
}
