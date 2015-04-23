package network

import (
	"db"
	"errorHandler"
	"time"
)

func HandleTaskStack(networkManager *NetworkManager) {
	duration := time.Duration(networkManager.Properties.ConnectionTaskStackSleepTime) * time.Millisecond

	for {
		rows, err := networkManager.Database.Query(db.SELECT([]string{"*", networkManager.Properties.ActionStack}))
		errorHandler.HandleError(err)

		defer rows.Close()
		for rows.Next() {
			var action Action

			err := rows.Scan(&action.ID, &action.Action, &action.Connection_condition, &action.Args)
			errorHandler.HandleError(err)

			action.Handle(networkManager)
		}
		errorHandler.HandleError(rows.Err())

		time.Sleep(duration)
	}
}
