package main

import (
	"database/sql"
	"db"
	"errorHandler"
	_ "github.com/go-sql-driver/mysql"
	"settings"
)

var (
	mysql_db   *sql.DB
	properties settings.Settings
)

func main() {
	setup()
}

func setup() {
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
