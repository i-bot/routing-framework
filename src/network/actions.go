package network

import (
	"os/exec"
	"strings"
	"strconv"
)

const (
	EXECUTE     = "EXEC"
	FORWARD_MSG = "FORWARD_MSG"
)

func HandleAction(action, action_description, msg string, networkManager *NetworkManager, identifier ConnectionIdentifier) {
	switch action {
	case EXECUTE:
		cmd := strings.Split(action_description, " ")[0]
		args := strings.Split(action_description, " ")[1:]
		args = append(args, msg, identifier.RemoteAddress.IP.String(), strconv.Itoa(identifier.LocalAddress.Port), strconv.Itoa(identifier.RemoteAddress.Port))
		command := exec.Command(cmd, args...)
		if networkManager.Properties.IsActionExecutionSynchronous {
			go command.Run()
		}else {
			command.Run()
		}
	case FORWARD_MSG:
		//Add message forwarding
	}
}
