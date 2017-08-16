package database

import (
	"database/sql"
	"log"
	"sync"

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
	db, err := sql.Open("sqlite3", "file:"+DB_LOCATION+DB_NAME+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}

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
	if _, err = db.Exec("PRAGMA temp_store=MEMORY;"); err != nil {
		log.Println("Error while setting temp_store=MEMORY:", err.Error())
	}

	return &DBConnector{db}, nil
}

func Reset() (e error) {
	resetMutex.Lock()
	defer resetMutex.Unlock()
	for _, table := range STS_TABLES {
		if _, err := conn.Exec(`DELETE FROM ` + table); err != nil {
			e = err
		}
	}
	return e
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
