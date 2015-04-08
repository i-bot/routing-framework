package main

import (
	"db"
	"fmt"
)

func main() {
	var db_connection db.DB_Connection
	mysql := db.MySQL{}
	db_connection = &mysql

	db_connection.Open("go", "go", "127.0.0.1", "3306", "godb")
	rows := db_connection.Query("SELECT * FROM animals")

	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		
		fmt.Print(id)
		fmt.Println("|" + name)
	}

	defer db_connection.Close()
}
