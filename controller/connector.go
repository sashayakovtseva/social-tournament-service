package controller

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	connectorSingleton *DBConnector
	connectorError     error
	connectorOnce      sync.Once
)

type DBConnector struct {
	*sql.DB
}

func newConnector() (*DBConnector, error) {
	db, err := sql.Open("sqlite3", "file:"+DB_LOCATION+DB_NAME+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	// create tables if they don't exist yet
	if _, err = db.Exec(CREATE_PLAYERS_TABLE); err != nil {
		db.Close()
		return nil, err
	}
	if _, err = db.Exec(CREATE_TOURNAMENTS_TABLE); err != nil {
		db.Close()
		return nil, err
	}
	if _, err = db.Exec(CREATE_P2T_TABLE); err != nil {
		db.Close()
		return nil, err
	}
	if _, err = db.Exec(CREATE_P2B_TABLE); err != nil {
		db.Close()
		return nil, err
	}

	// define pragmas for performance increase
	if _, err = db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		log.Println("Error while setting journal_mode=WAL:", err.Error())
	}
	if _, err = db.Exec("PRAGMA synchronous=OFF;"); err != nil {
		log.Println("Error while setting synchronous=OFF:", err.Error())
	}
	if _, err = db.Exec("PRAGMA temp_store=MEMORY;"); err != nil {
		log.Println("Error while setting temp_store=MEMORY:", err.Error())
	}
	if _, err = db.Exec("PRAGMA locking_mode=EXCLUSIVE;"); err != nil {
		log.Println("Error while setting locking_mode=EXCLUSIVE:", err.Error())
	}


	return &DBConnector{db}, nil
}

func (conn *DBConnector) Reset() (e error) {
	for _, table := range STS_TABLES {
		if _, err := conn.Exec(`DELETE FROM ` + table); err != nil {
			e = err
		}
	}
	return e
}

func GetConnector() (*DBConnector, error) {
	connectorOnce.Do(func() {
		connectorSingleton, connectorError = newConnector()

	})
	return connectorSingleton, connectorError
}
