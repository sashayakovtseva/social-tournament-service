package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	TournamentConn             *TournamentConnector
	ErrTournamentAlreadyExists = errors.New("tournament already exists")
)

type TournamentConnector struct {
	insert           *sql.Stmt
	update           *sql.Stmt
	read             *sql.Stmt
	join             *sql.Stmt
	insertBackPlayer *sql.Stmt
	readParticipants *sql.Stmt
	readBackers      *sql.Stmt
	statements       []*sql.Stmt
}

func newTournamentConnector() (*TournamentConnector, error) {
	var err error

	tournamentConn := new(TournamentConnector)

	tournamentConn.insert, err = prepareAndAdd(tournamentConn.statements, `INSERT INTO %s(%s,%s) values(?,?)`,
		tournamentsTableName, tournamentIDColName, depositColName)
	if haveToFailAndClose(tournamentConn, err) {
		return nil, err
	}

	tournamentConn.update, err = prepareAndAdd(tournamentConn.statements, `UPDATE %s SET %s = ? WHERE %s = ?`,
		tournamentsTableName, finishedColName, tournamentIDColName)
	if haveToFailAndClose(tournamentConn, err) {
		return nil, err
	}

	tournamentConn.read, err = prepareAndAdd(tournamentConn.statements,
		`SELECT %s, %s, %s FROM %s WHERE %s = ?`,
		tournamentIDColName, depositColName, finishedColName, tournamentsTableName, tournamentIDColName)
	if haveToFailAndClose(tournamentConn, err) {
		return nil, err
	}

	tournamentConn.join, err = prepareAndAdd(tournamentConn.statements, `INSERT INTO %s(%s,%s) values(?,?)`,
		p2tTableName, playerIDColName, tournamentIDColName)
	if haveToFailAndClose(tournamentConn, err) {
		return nil, err
	}

	tournamentConn.insertBackPlayer, err = prepareAndAdd(tournamentConn.statements,
		`INSERT INTO %s(%s,%s, %s) values(?,?, ?)`,
		p2bTableName, playerIDColName, tournamentIDColName, backerColName)
	if haveToFailAndClose(tournamentConn, err) {
		return nil, err
	}

	tournamentConn.readParticipants, err = prepareAndAdd(tournamentConn.statements,
		`SELECT %s, %s FROM %s INNER JOIN %s  USING (%s) WHERE %s = ?`,
		playerIDColName, balanceColName, playersTableName, p2tTableName,
		playerIDColName, tournamentIDColName)
	if haveToFailAndClose(tournamentConn, err) {
		return nil, err
	}

	tournamentConn.readBackers, err = prepareAndAdd(tournamentConn.statements,
		`SELECT %s.%s, %s FROM %s INNER JOIN %s ON %s.%s  = %s.%s WHERE %s = ? AND %s.%s = ?`,
		playersTableName, playerIDColName, balanceColName, playersTableName, p2bTableName, p2bTableName,
		backerColName, playersTableName, playerIDColName, tournamentIDColName, p2bTableName, playerIDColName)
	if haveToFailAndClose(tournamentConn, err) {
		return nil, err
	}

	return tournamentConn, nil
}

func (connector *TournamentConnector) Close() {
	for _, stmt := range connector.statements {
		checkAndClose(stmt)
	}
}

func (connector *TournamentConnector) Create(tournament *entity.Tournament) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	_, err := connector.insert.Exec(tournament.ID(), tournament.Deposit())
	return replaceConstraintWithCustom(err, ErrTournamentAlreadyExists)
}

func (connector *TournamentConnector) Read(id string) *entity.Tournament {
	var tournamentID string
	var tournamentDeposit float32
	var tournamentIsFinished bool

	resetMutex.RLock()
	defer resetMutex.RUnlock()

	err := connector.read.QueryRow(id).Scan(&tournamentID, &tournamentDeposit, &tournamentIsFinished)
	if err != nil {
		return nil
	}
	return entity.NewTournament(tournamentID, tournamentDeposit, tournamentIsFinished)
}

func (connector *TournamentConnector) UpdateWithTx(tx *sql.Tx, tournament *entity.Tournament) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	preparedUpdateTx := tx.Stmt(connector.update)
	_, err := preparedUpdateTx.Exec(tournament.IsFinished(), tournament.ID())
	return err
}

func (connector *TournamentConnector) UpdateResult(tournament *entity.Tournament, players []*entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	err = PlayerConn.UpdateWithTx(tx, players...)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = connector.UpdateWithTx(tx, tournament)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (connector *TournamentConnector) SelectBackPlayers(tournamentID, playerID string) ([]*entity.Player, error) {
	var backers []*entity.Player
	var backerID string
	var balance float32

	resetMutex.RLock()
	defer resetMutex.RUnlock()

	rows, err := connector.readBackers.Query(tournamentID, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&backerID, &balance)
		backers = append(backers, entity.NewPlayer(backerID, balance))
	}
	return backers, rows.Err()
}

func (connector *TournamentConnector) SelectParticipants(tournamentID string) ([]*entity.Player, error) {
	var participants []*entity.Player
	var playerID string
	var balance float32

	resetMutex.RLock()
	defer resetMutex.RUnlock()

	rows, err := connector.readParticipants.Query(tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&playerID, &balance)
		participants = append(participants, entity.NewPlayer(playerID, balance))
	}
	return participants, rows.Err()
}

func (connector *TournamentConnector) UpdateJoin(tournament *entity.Tournament,
	player *entity.Player, backers []*entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	PlayerConn.UpdateWithTx(tx, player)
	if err != nil {
		tx.Rollback()
		return err
	}
	PlayerConn.UpdateWithTx(tx, backers...)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Stmt(connector.join).Exec(player.ID(), tournament.ID())
	if err != nil {
		tx.Rollback()
		return err
	}
	txPreparedBackPlayer := tx.Stmt(connector.insertBackPlayer)
	for _, backer := range backers {
		_, err := txPreparedBackPlayer.Exec(player.ID(), tournament.ID(), backer.ID())
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func init() {
	var err error
	TournamentConn, err = newTournamentConnector()
	if err != nil {
		log.Fatal(err)
	}
}
