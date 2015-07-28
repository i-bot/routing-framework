package network

import (
	"errorHandler"
	"strings"

	"github.com/i-bot/mysqlParser"
)

func HandleConnect(networkManager *NetworkManager, identifier string) {
	ip, localport, remoteport := networkManager.ConvertToStrings(identifier)
	table := networkManager.Properties.Connections

	_, err := networkManager.Database.Exec(mysqlParser.INSERT_INTO([]string{
		table, "ip,localport,remoteport", "'" + ip + "'," + localport + "," + remoteport}))
	errorHandler.HandleError(err)

	scanAndHandleRows(networkManager.Properties.OnOpen, "", networkManager, identifier, "true")
}

func HandleRead(msg string, networkManager *NetworkManager, identifier string) {
	table := networkManager.Properties.OnRead

	scanAndHandleRows(table, msg, networkManager, identifier, "\""+msg+"\" REGEXP "+table+".msg_regex")
}

func HandleWrite(msg string, networkManager *NetworkManager, identifier string) {
	table := networkManager.Properties.OnWrite

	scanAndHandleRows(table, msg, networkManager, identifier, "\""+msg+"\" REGEXP "+table+".msg_regex")
}

func HandleClose(networkManager *NetworkManager, identifier string) {
	ip, localport, remoteport := networkManager.ConvertToStrings(identifier)
	table := networkManager.Properties.Connections

	_, err := networkManager.Database.Exec(mysqlParser.DELETE([]string{
		table, table + ".ip=\"" + ip + "\" AND " + table +
			".localport=" + localport + " AND " + table + ".remoteport=" + remoteport}))
	errorHandler.HandleError(err)

	scanAndHandleRows(networkManager.Properties.OnClose, "", networkManager, identifier, "true")
}

func scanAndHandleRows(table string, msg string, networkManager *NetworkManager, identifier string, actionCondition string) {
	ip, localport, remoteport := networkManager.ConvertToStrings(identifier)

	rows, err := networkManager.Database.Query(mysqlParser.SELECT([]string{"id, connectionCondition, action, args", table, actionCondition}))
	errorHandler.HandleError(err)

	defer rows.Close()
	for rows.Next() {
		var action Action

		err := rows.Scan(&action.ID, &action.ConnectionCondition, &action.Action, &action.Args)
		errorHandler.HandleError(err)

		filteredConnections := "filteredConnections"
		matchingConnections, err := networkManager.Database.Query(
			mysqlParser.SELECT([]string{
				"*",
				strings.TrimSuffix(
					mysqlParser.AS([]string{
						mysqlParser.SELECT([]string{
							"*",
							networkManager.Properties.Connections,
							action.ConnectionCondition}),
						filteredConnections}),
					";"),
				filteredConnections + ".ip=\"" + ip + "\" AND " + filteredConnections + ".localport=" + localport + " AND " + filteredConnections + ".remoteport=" + remoteport}))
		errorHandler.HandleError(err)

		size := 0

		for ; matchingConnections.Next(); size++ {
		}

		if size == 1 {
			action.Msg = msg
			action.Identifier = identifier
			handleAction(action, networkManager)
		}

		errorHandler.HandleError(matchingConnections.Err())
		matchingConnections.Close()
	}
	errorHandler.HandleError(rows.Err())
}

func handleAction(action Action, networkManager *NetworkManager) {
	switch action.Action {
	case CLOSE:
		networkManager.Close(action.Identifier)

	case WRITE:
		networkManager.Write(action.Identifier, action.Args)

	default:
		action.Handle(networkManager)
	}
}
