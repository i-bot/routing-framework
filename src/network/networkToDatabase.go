package network

import (
	"db"
	"database/sql"
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

func scanAndHandleRows(rows *sql.Rows, msg string, networkManager *NetworkManager, identifier ConnectionIdentifier, ip, localport, remoteport string){
	defer rows.Close()
	for rows.Next() {
		var id int
		var connection_condition, action, action_description string
		err := rows.Scan(&id, &connection_condition, &action, &action_description)
		errorHandler.HandleError(err)

		matchingConnections, err := networkManager.Database.Query(db.SELECT([]string{
			"*", "(" + db.SELECT([]string{"*", networkManager.Properties.Connections, connection_condition}) + ")",
			networkManager.Properties.Connections + ".ip=" + ip + "AND " + networkManager.Properties.Connections + ".localport=" + localport +
				"AND " + networkManager.Properties.Connections + ".remoteport=" + remoteport}))

		size := 0

		for ; rows.Next(); size++ {
		}

		if size == 1 {
			HandleAction(action, action_description, msg, networkManager, identifier)
		}

		matchingConnections.Close()
		errorHandler.HandleError(matchingConnections.Err())
	}
	errorHandler.HandleError(rows.Err())
}
