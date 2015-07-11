package framework

import (
	"database/sql"
	"db"
	"errorHandler"
	_ "github.com/go-sql-driver/mysql"
	"network"
	"settings"
)

var (
	mysql_db       *sql.DB
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

	mysql_db, err = sql.Open("mysql", db.OPEN([]string{properties.Username, properties.Password, properties.IP, properties.Port, properties.DB_Name}))
	errorHandler.HandleError(err)

	for _, element := range properties.Databases {
		_, err = mysql_db.Exec(db.DROP_TABLE(element[:1]))
		errorHandler.HandleError(err)
		_, err = mysql_db.Exec(db.CREATE_TABLE(element))
		errorHandler.HandleError(err)
	}

	for _, element := range properties.Values {
		_, err = mysql_db.Exec(db.INSERT_INTO(element))
		errorHandler.HandleError(err)
	}

	return
}

func setupNetworkManager() {
	networkManager = &network.NetworkManager{mysql_db, properties}
	networkManager.Init()
}

func startHandlingTaskQueue() {
	network.HandleTaskQueue(networkManager)
}
