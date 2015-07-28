package framework

import (
	"database/sql"
	"errorHandler"
	"network"
	"settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/i-bot/mysqlParser"
)

var (
	mysqlDB        *sql.DB
	properties     *settings.Settings
	networkManager *network.NetworkManager
)

func Start(props *settings.Settings) {
	properties = props

	setupDatabase()
	setupNetworkManager()
	startHandlingTaskQueue()
}

func setupDatabase() {
	var err error

	mysqlDB, err = sql.Open("mysql", mysqlParser.OPEN([]string{properties.Username, properties.Password, properties.IP, properties.Port, properties.DBName}))
	errorHandler.HandleError(err)

	for _, element := range properties.Databases {
		_, err = mysqlDB.Exec(mysqlParser.DROP_TABLE(element[:1]))
		errorHandler.HandleError(err)
		_, err = mysqlDB.Exec(mysqlParser.CREATE_TABLE(element))
		errorHandler.HandleError(err)
	}

	for _, element := range properties.Values {
		_, err = mysqlDB.Exec(mysqlParser.INSERT_INTO(element))
		errorHandler.HandleError(err)
	}

	return
}

func setupNetworkManager() {
	networkManager = &network.NetworkManager{Database: mysqlDB, Properties: properties}
	networkManager.Init()
}

func startHandlingTaskQueue() {
	network.HandleTaskQueue(networkManager)
}
