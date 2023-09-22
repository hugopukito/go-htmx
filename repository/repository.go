package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var err error

func InitDB() {
	DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/htmx?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
}
