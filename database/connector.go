package database

import (
	"database/sql"
	"log"
	"sync"

	// imported for driver registration
	_ "github.com/mattn/go-sqlite3"
)

var (
	conn       *DBConnector
	resetMutex sync.RWMutex
)

type DBConnector struct {
	*sql.DB
}

func newConnector() (*DBConnector, error) {
	db, err := sql.Open("sqlite3", "file:"+dbLocation+dbName+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}
	for _, SQL := range SQLs {
		if err = execSQLCloseOnFail(db, SQL); err != nil {
			return nil, err
		}
	}
	return &DBConnector{db}, nil
}

func execSQLCloseOnFail(db *sql.DB, SQL string) (err error) {
	if _, err = db.Exec(SQL); err != nil {
		db.Close()
	}
	return
}

func Reset() (e error) {
	resetMutex.Lock()
	defer resetMutex.Unlock()

	for _, table := range stsTables {
		if _, err := conn.Exec(`DELETE FROM ` + table); err != nil {
			e = err
		}
	}
	return
}

func Close() {
	conn.Close()
	PlayerConn.Close()
	TournamentConn.Close()
}

func init() {
	var err error
	conn, err = newConnector()
	if err != nil {
		log.Fatal(err)
	}
}
