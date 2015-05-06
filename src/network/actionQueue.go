package network

import (
	"db"
	"errorHandler"
	"time"
)

func HandleTaskQueue(networkManager *NetworkManager) {
	duration := time.Duration(networkManager.Properties.ConnectionTaskStackSleepTime) * time.Millisecond

	for {
		rows, err := networkManager.Database.Query(db.SELECT([]string{"*", networkManager.Properties.ActionQueue}))
		errorHandler.HandleError(err)

		defer rows.Close()
		for rows.Next() {
			var action Action

			err := rows.Scan(&action.ID, &action.Action, &action.Connection_condition, &action.Args)
			errorHandler.HandleError(err)

			action.Handle(networkManager)
		}
		errorHandler.HandleError(rows.Err())
		
		_, err = networkManager.Database.Query(db.UPDATE_PRIMARY_KEY([]string{networkManager.Properties.ActionQueue}))
		errorHandler.HandleError(err)

		time.Sleep(duration)
	}
}
