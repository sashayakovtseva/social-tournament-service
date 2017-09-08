package database

import "fmt"

const (
	dbLocation = "./"
	dbName     = "sts.db"

	playersTableName     = "players"
	tournamentsTableName = "tournaments"
	p2tTableName         = "p2t"
	p2bTableName         = "p2b"

	playerIDColName     = "pid"
	tournamentIDColName = "tid"
	balanceColName      = "balance"
	depositColName      = "deposit"
	finishedColName     = "finished"
	backerColName       = "bid"
)

var (
	createPlayersTableIfNotExists = fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s(%s TEXT PRIMARY KEY CHECK(%s <> ""),
									   %s INTEGER NOT NULL CHECK (%s >= 0))`,
		playersTableName, playerIDColName, playerIDColName, balanceColName, balanceColName)

	createTournamentsTableIfNotExists = fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s(%s TEXT PRIMARY KEY CHECK(%s <> ""),
									   %s INTEGER NOT NULL CHECK (%s >= 0),
									   %s BOOL NOT NULL DEFAULT 0)`,
		tournamentsTableName, tournamentIDColName, tournamentIDColName, depositColName, depositColName, finishedColName)

	createP2tTableIfNotExists = fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s(%s TEXT REFERENCES %s(%s) ON DELETE CASCADE,
									   %s TEXT REFERENCES %s(%s) ON DELETE CASCADE,
									   PRIMARY KEY (%s,%s))`,
		p2tTableName, playerIDColName, playersTableName, playerIDColName,
		tournamentIDColName, tournamentsTableName, tournamentIDColName,
		playerIDColName, tournamentIDColName)

	createP2bTableIfNotExists = fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s(%s TEXT REFERENCES %s(%s) ON DELETE CASCADE,
									   %s TEXT REFERENCES %s(%s) ON DELETE CASCADE,
									   %s TEXT REFERENCES %s(%s) ON DELETE CASCADE,
									   PRIMARY KEY (%s,%s, %s),
									   FOREIGN KEY (%s, %s) REFERENCES %s(%s, %s))`,
		p2bTableName, playerIDColName, playersTableName, playerIDColName,
		tournamentIDColName, tournamentsTableName, tournamentIDColName,
		backerColName, playersTableName, playerIDColName,
		playerIDColName, tournamentIDColName, backerColName,
		playerIDColName, tournamentIDColName, p2tTableName, playerIDColName, tournamentIDColName)

	pragmaJournalMode = "PRAGMA journal_mode=WAL"
	pragmaTempStore   = "PRAGMA temp_store=MEMORY"

	SQLs = []string{createP2bTableIfNotExists, createP2tTableIfNotExists,
		createPlayersTableIfNotExists, createTournamentsTableIfNotExists,
		pragmaJournalMode, pragmaTempStore}
	stsTables = []string{p2bTableName, p2tTableName, tournamentsTableName, playersTableName}
)
