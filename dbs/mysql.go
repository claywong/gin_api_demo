package dbs

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Conns *sql.DB

func init() {
	var err error
	dns := "root:123456@tcp(localhost:3306)/ppgo_api_demo_gin?parseTime=true"
	Conns, err = sql.Open("mysql", dns)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = Conns.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	Conns.SetMaxIdleConns(20)
	Conns.SetMaxOpenConns(20)
}
