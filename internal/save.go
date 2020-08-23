package internal

import (
	"database/sql"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"

	"github.com/OakAnderson/internetTester/nettest"
)

// CSV implements nettest.Saver to save results into a file
type CSV struct {
	file string
}

// Exec saves the args into the csv file
func (csv CSV) Exec(args ...interface{}) (sql.Result, error) {
	var row string
	for _, arg := range args[:len(args)-1] {
		row += fmt.Sprintf("%v,", arg)
	}
	row += fmt.Sprintf("%v\n", args[len(args)-1])

	file, err := os.OpenFile(csv.file, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(row)

	file.Close()

	return nil, err
}

// connDatabase connect to mysql database and return it
func connDatabase(user, password, database string) (db *sql.DB, err error) {
	return sql.Open("mysql", user+":"+password+"@/"+database)
}

func createTableIfNotExists(db *sql.DB) error {
	dbTable, err := ioutil.ReadFile(build.Default.GOPATH + "/src/github.com/OakAnderson/internetTester/database/db.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(dbTable))

	return err
}

// SaveMysql return the sql.Stmt for insert to database
func SaveMysql(user, password, database string) (nettest.Saver, error) {
	db, err := connDatabase(user, password, database)
	if err != nil {
		return nil, err
	}

	err = createTableIfNotExists(db)
	if err != nil {
		return nil, err
	}

	return db.Prepare(
		"INSERT speedtest SET dt=?,latency=?,jitter=?,download=?,upload=?,packetLoss=?,hardware=?,serverId=?,ip=?,name=?,location=?,host=?",
	)
}

// SaveCSV return a nettest.Saver to save into a file
func SaveCSV(filename, columns string) (nettest.Saver, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() == 0 {
		_, err = file.Write(
			[]byte(columns))

		if err != nil {
			return nil, err
		}
	}
	err = file.Close()

	return CSV{filename}, err
}
