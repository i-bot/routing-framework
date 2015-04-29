package network

import (
	"db"
	"errorHandler"
)

func HandleOpen(networkManager *NetworkManager, identifier string) {
	ip, localport, remoteport := networkManager.ConvertToStrings(identifier)

	_, err := networkManager.Database.Query(db.INSERT_INTO([]string{
		networkManager.Properties.Connections, "ip,localport,remoteport", "'" + ip + "'," + localport + "," + remoteport}))
	errorHandler.HandleError(err)

	scanAndHandleRows(networkManager.Properties.OnOpenActions, "", networkManager, identifier)
}

func HandleRead(msg string, networkManager *NetworkManager, identifier string) {
	scanAndHandleRows(networkManager.Properties.OnReadActions, msg, networkManager, identifier)
}

func HandleWrite(msg string, networkManager *NetworkManager, identifier string) {
	scanAndHandleRows(networkManager.Properties.OnWriteActions, msg, networkManager, identifier)
}

func HandleClose(networkManager *NetworkManager, identifier string) {
	ip, localport, remoteport := networkManager.ConvertToStrings(identifier)

	_, err := networkManager.Database.Query(db.DELETE([]string{
		networkManager.Properties.Connections, networkManager.Properties.Connections + ".ip=" + ip + "AND " + networkManager.Properties.Connections +
			".localport=" + localport + "AND " + networkManager.Properties.Connections + ".remoteport=" + remoteport}))
	errorHandler.HandleError(err)

	scanAndHandleRows(networkManager.Properties.OnCloseActions, "", networkManager, identifier)
}

func scanAndHandleRows(table string, msg string, networkManager *NetworkManager, identifier string) {
	ip, localport, remoteport := networkManager.ConvertToStrings(identifier)

	rows, err := networkManager.Database.Query(db.SELECT([]string{"*", table}))
	errorHandler.HandleError(err)

	defer rows.Close()
	for rows.Next() {
		var action Action

		err := rows.Scan(&action.ID, &action.Action, &action.Connection_condition, &action.Args)
		errorHandler.HandleError(err)

		matchingConnections, err := networkManager.Database.Query(db.SELECT([]string{
			"*", "(" + db.SELECT([]string{"*", networkManager.Properties.Connections, action.Connection_condition}) + ")",
			networkManager.Properties.Connections + ".ip=" + ip + "AND " + networkManager.Properties.Connections + ".localport=" + localport +
				"AND " + networkManager.Properties.Connections + ".remoteport=" + remoteport}))

		size := 0

		for ; rows.Next(); size++ {
		}

		if size == 1 {
			action.Msg = msg
			action.Identifier = identifier
			action.Handle(networkManager)
		}

		errorHandler.HandleError(matchingConnections.Err())
		matchingConnections.Close()
	}
	errorHandler.HandleError(rows.Err())
}
