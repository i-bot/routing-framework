package network

import (
	"errorHandler"
	"strconv"
	"time"

	"github.com/i-bot/mysqlParser"
)

func HandleTaskQueue(networkManager *NetworkManager) {
	duration := time.Duration(networkManager.Properties.ActionQueueSleepTime) * time.Millisecond

	for {
		rows, err := networkManager.Database.Query(mysqlParser.SELECT("*", networkManager.Properties.ActionQueue))
		errorHandler.HandleError(err)

		for rows.Next() {
			var action Action

			err := rows.Scan(&action.ID, &action.ConnectionCondition, &action.Action, &action.Args)
			errorHandler.HandleError(err)

			action.Handle(networkManager)

			_, err = networkManager.Database.Exec(mysqlParser.DELETE(networkManager.Properties.ActionQueue, "actionQueue.id="+strconv.Itoa(action.ID)))
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

	_, err = tx.Prepare(mysqlParser.SET("@count=0"))
	errorHandler.HandleError(err)

	_, err = tx.Prepare(mysqlParser.UPDATE(table, table+".id = @count:= @count + 1"))
	errorHandler.HandleError(err)

	err = tx.Commit()
	errorHandler.HandleError(err)
}
