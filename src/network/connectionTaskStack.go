package network

import (
	"db"
	"errorHandler"
	"strconv"
	"strings"
	"time"
)

var (
	CONNECT = "CONNECT"
	LISTEN  = "LISTEN"
	CLOSE   = "CLOSE"
	WRITE   = "WRITE"
)

type ITask interface {
	start(networkManager *NetworkManager)
}

type Task struct {
	id                               int
	task, connection_condition, args string
}

func HandleTaskStack(networkManager *NetworkManager) {
	duration := time.Duration(networkManager.Properties.ConnectionTaskStackSleepTime) * time.Millisecond

	for {
		rows, err := networkManager.Database.Query(db.SELECT([]string{"*", networkManager.Properties.ConnectionsTaskStack}))
		errorHandler.HandleError(err)

		defer rows.Close()
		for rows.Next() {
			var task Task
			err := rows.Scan(&task.id, &task.task, &task.connection_condition, &task.args)
			errorHandler.HandleError(err)

			task.start(networkManager)
		}
		errorHandler.HandleError(rows.Err())

		time.Sleep(duration)
	}
}

func (task Task) start(networkManager *NetworkManager) {
	switch task.task {
	case CONNECT:
		split := strings.Split(task.args, ":")
		remoteport, err := strconv.Atoi(split[1])
		errorHandler.HandleError(err)

		networkManager.Connect(split[0], remoteport)
	case LISTEN:
		split := strings.Split(task.args, ":")
		localport, err := strconv.Atoi(split[0])
		errorHandler.HandleError(err)

		networkManager.Listen(localport)
	case CLOSE:
		identifiers := getMatchingConnections(&task, networkManager)

		for _, identifier := range identifiers {
			networkManager.Close(identifier)
		}
	case WRITE:
		identifiers := getMatchingConnections(&task, networkManager)

		for _, identifier := range identifiers {
			networkManager.Write(identifier, task.args)
		}
	}
}

func getMatchingConnections(task *Task, networkManager *NetworkManager) (identifiers []ConnectionIdentifier) {
	rows, err := networkManager.Database.Query(db.SELECT([]string{"*", networkManager.Properties.Connections, task.connection_condition}))
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
