package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type DB_Connection interface {
	Open(user, password, ip, port, db_name string)
	Ping()
	Query(request string) (*sql.Rows)
	Close()
}

type MySQL struct {
	db *sql.DB
}

func (mySQL *MySQL) Open(user, password, ip, port, db_name string) {
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+ip+":"+port+")/"+db_name)
	mySQL.db = db

	handleError(err)
}

func (mySQL *MySQL) Ping() {
	err := mySQL.db.Ping()

	handleError(err)
}

func (mySQL *MySQL) Query(request string) (*sql.Rows) {
	rows, err := mySQL.db.Query(request)

	handleError(err)
	
	return rows
}

func (mySQL *MySQL) Close() {
	err := mySQL.db.Close()

	handleError(err)
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e.Error)
		os.Exit(1)
	}
}
