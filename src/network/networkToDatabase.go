package network

import (
	"database/sql"
	"db"
	"errorHandler"
	"strconv"
)

func AddConnection(networkManager *NetworkManager, identifier ConnectionIdentifier) {
	ip := identifier.RemoteAddress.IP.String()
	localport := strconv.Itoa(identifier.LocalAddress.Port)
	remoteport := strconv.Itoa(identifier.RemoteAddress.Port)

	_, err := networkManager.Database.Query(db.INSERT_INTO([]string{
		networkManager.Properties.Connections, "ip,localport,remoteport", "'" + ip + "'," + localport + "," + remoteport}))
	errorHandler.HandleError(err)

	rows, err := networkManager.Database.Query(db.SELECT([]string{"*", networkManager.Properties.OnOpenActions}))
	errorHandler.HandleError(err)

	scanAndHandleRows(rows, "", networkManager, identifier, ip, localport, remoteport)
}

func scanAndHandleRows(rows *sql.Rows, msg string, networkManager *NetworkManager, identifier ConnectionIdentifier, ip, localport, remoteport string) {
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
