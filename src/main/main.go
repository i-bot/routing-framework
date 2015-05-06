package main

import (
	"database/sql"
	"db"
	"errorHandler"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"network"
	"os"
	"settings"
)

var (
	mysql_db       *sql.DB
	properties     *settings.Settings
	networkManager *network.NetworkManager
)

func main() {
	loadSettings(os.Args)
	setupDatabase()
	setupNetworkManager()
	startHandlingTaskQueue()
}

func loadSettings(args []string) {
	if len(args) == 1 {
		properties = settings.LoadSettings(args[0])
	} else {
		errorHandler.HandleError(errors.New("Pass the location of the configuration file"))
	}
}

func setupDatabase() {
	var err error

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

func setupNetworkManager() {
	networkManager = &network.NetworkManager{mysql_db, properties}
}

func startHandlingTaskQueue() {
	go network.HandleTaskQueue(networkManager)
}
