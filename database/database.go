package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Mysql Driver
)

// ConnDatabase is
func ConnDatabase() (db *sql.DB, err error) {
	user, dbname := "oak", "internetTester"
	db, err = sql.Open("mysql", user+":@/"+dbname)
	return
}
