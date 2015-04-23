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
	ID                                               int
	Action, Connection_condition, Args string
}

func HandleAction(action, action_description, msg string, networkManager *NetworkManager, identifier ConnectionIdentifier) {
	switch action {
	case EXECUTE:
		cmd := strings.Split(action_description, " ")[0]
		args := strings.Split(action_description, " ")[1:]
		args = append(args, msg, identifier.RemoteAddress.IP.String(), strconv.Itoa(identifier.LocalAddress.Port), strconv.Itoa(identifier.RemoteAddress.Port))
		command := exec.Command(cmd, args...)
		if networkManager.Properties.IsActionExecutionSynchronous {
			go command.Run()
		} else {
			command.Run()
		}
	case FORWARD_MSG:
		//Add message forwarding
	}
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
		identifiers := getMatchingConnections(&action, networkManager)

		for _, identifier := range identifiers {
			networkManager.Close(identifier)
		}
	case WRITE:
		identifiers := getMatchingConnections(&action, networkManager)

		for _, identifier := range identifiers {
			networkManager.Write(identifier, action.Args)
		}
	}
}

func getMatchingConnections(action *Action, networkManager *NetworkManager) (identifiers []ConnectionIdentifier) {
	rows, err := networkManager.Database.Query(db.SELECT([]string{"*", networkManager.Properties.Connections, action.Connection_condition}))
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
