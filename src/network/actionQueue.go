package network

import (
	"db"
	"errorHandler"
	"strconv"
	"time"
)

func HandleTaskQueue(networkManager *NetworkManager) {
	duration := time.Duration(networkManager.Properties.ActionQueueSleepTime) * time.Millisecond

	for {
		rows, err := networkManager.Database.Query(db.SELECT([]string{"*", networkManager.Properties.ActionQueue}))
		errorHandler.HandleError(err)

		for rows.Next() {
			var action Action

			err := rows.Scan(&action.ID, &action.Connection_condition, &action.Action, &action.Args)
			errorHandler.HandleError(err)

			action.Handle(networkManager)

			_, err = networkManager.Database.Exec(db.DELETE([]string{networkManager.Properties.ActionQueue, "actionQueue.id=" + strconv.Itoa(action.ID)}))
			errorHandler.HandleError(err)
		}
		errorHandler.HandleError(rows.Err())

		updateID(networkManager, networkManager.Properties.ActionQueue)
		errorHandler.HandleError(err)

		rows.Close()
		time.Sleep(duration)
	}
}

func updateID(networkManager *NetworkManager, table string) {
	tx, err := networkManager.Database.Begin()
	errorHandler.HandleError(err)

	defer tx.Rollback()

	_, err = tx.Prepare("SET @count=0;")
	errorHandler.HandleError(err)

	_, err = tx.Prepare("UPDATE " + table + " SET " + table + ".id = @count:= @count + 1;")
	errorHandler.HandleError(err)

	err = tx.Commit()
	errorHandler.HandleError(err)
}
