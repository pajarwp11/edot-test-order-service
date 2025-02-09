package mysql

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	MySQL *sqlx.DB
)

func Connect() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/test_database"
	MySQL, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
