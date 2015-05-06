package main

import (
	"database/sql"
	"db"
	"errorHandler"
	_ "github.com/go-sql-driver/mysql"
	"settings"
	"network"
)

var (
	mysql_db   *sql.DB
	properties *settings.Settings
	networkManager *network.NetworkManager
)

func main() {
	setupDatabase()
	setupNetworkManager()
	startHandlingTaskQueue()
}

func setupDatabase() {
	var err error

	properties = settings.LoadSettings()

	mysql_db, err = sql.Open("mysql", db.OPEN([]string{properties.Username, properties.Password, properties.IP, properties.Port, properties.DB_Name}))
	errorHandler.HandleError(err)

	for _, element := range properties.Databases {
		_, err = mysql_db.Query(db.DROP_TABLE(element[:1]))
		errorHandler.HandleError(err)
		_, err = mysql_db.Query(db.CREATE_TABLE(element))
		errorHandler.HandleError(err)
	}

	return
}

func setupNetworkManager(){
	networkManager = &network.NetworkManager{mysql_db, properties}
}

func startHandlingTaskQueue(){
	go network.HandleTaskQueue(networkManager)
}
