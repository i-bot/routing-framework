package network

import (
	"db"
	"errorHandler"
	"os/exec"
	"strconv"
	"strings"
)

const (
	CONNECT     = "CONNECT"
	LISTEN      = "LISTEN"
	CLOSE       = "CLOSE"
	WRITE       = "WRITE"
	EXECUTE     = "EXEC"
	FORWARD_MSG = "FORWARD_MSG"
)

type IAction interface {
	Handle(networkManager *NetworkManager)
}

type Action struct {
	ID                                 int
	Action, Connection_condition, Args string

	Msg        string
	Identifier ConnectionIdentifier
}

func (action Action) Handle(networkManager *NetworkManager) {
	switch action.Action {
	case CONNECT:
		split := strings.Split(action.Args, ":")
		remoteport, err := strconv.Atoi(split[1])
		errorHandler.HandleError(err)

		networkManager.Connect(split[0], remoteport)

	case LISTEN:
		split := strings.Split(action.Args, ":")
		localport, err := strconv.Atoi(split[0])
		errorHandler.HandleError(err)

		networkManager.Listen(localport)

	case CLOSE:
		identifiers := getMatchingConnections(action.Connection_condition, networkManager)

		for _, identifier := range identifiers {
			networkManager.Close(identifier)
		}

	case WRITE:
		identifiers := getMatchingConnections(action.Connection_condition, networkManager)

		for _, identifier := range identifiers {
			networkManager.Write(identifier, action.Args)
		}

	case EXECUTE:
		cmd := strings.Split(action.Args, " ")[0]
		args := strings.Split(action.Args, " ")[1:]
		args = append(args, action.Msg, action.Identifier.RemoteAddress.IP.String(),
			strconv.Itoa(action.Identifier.LocalAddress.Port), strconv.Itoa(action.Identifier.RemoteAddress.Port))

		command := exec.Command(cmd, args...)

		if !networkManager.Properties.IsActionExecutionSynchronous {
			go command.Run()
		} else {
			command.Run()
		}

	case FORWARD_MSG:
		identifiers := getMatchingConnections(action.Args, networkManager)

		for _, identifier := range identifiers {
			networkManager.Write(identifier, action.Args)
		}
	}
}

func getMatchingConnections(where string, networkManager *NetworkManager) (identifiers []ConnectionIdentifier) {
	rows, err := networkManager.Database.Query(db.SELECT([]string{"*", networkManager.Properties.Connections, where}))
	errorHandler.HandleError(err)

	defer rows.Close()
	for rows.Next() {
		var ip string
		var localport, remoteport int
		err := rows.Scan(&ip, &localport, &remoteport)
		errorHandler.HandleError(err)

		identifiers = append(identifiers, networkManager.ConvertToConnectionIdentifier(ip, localport, remoteport))
	}
	errorHandler.HandleError(rows.Err())

	return
}
