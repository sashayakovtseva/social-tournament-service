package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	TournamentConn             *TournamentConnector
	ErrTournamentAlreadyExists = errors.New("tournament already exists")
)

type TournamentConnector struct {
	insert                     *sql.Stmt
	preparedUpdateTournament   *sql.Stmt
	preparedSelectTournament   *sql.Stmt
	preparedJoinTournament     *sql.Stmt
	preparedInsertBackPlayer   *sql.Stmt
	preparedSelectParticipants *sql.Stmt
	preparedSelectBackers      *sql.Stmt
}

func newTournamentConnector() (*TournamentConnector, error) {
	var err error
	tournamentConn := new(TournamentConnector)
	tournamentConn.insert, err = conn.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		tournamentsTableName, tournamentIDColName, depositColName))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedUpdateTournament, err = conn.Prepare(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?`,
		tournamentsTableName, finishedColName, tournamentIDColName))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedSelectTournament, err = conn.Prepare(fmt.Sprintf(
		`SELECT %s, %s, %s FROM %s WHERE %s = ?`,
		tournamentIDColName, depositColName, finishedColName, tournamentsTableName, tournamentIDColName))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedJoinTournament, err = conn.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		p2tTableName, playerIDColName, tournamentIDColName))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedInsertBackPlayer, err = conn.Prepare(fmt.Sprintf(
		`INSERT INTO %s(%s,%s, %s) values(?,?, ?)`,
		p2bTableName, playerIDColName, tournamentIDColName, backerColName))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedSelectParticipants, err = conn.Prepare(fmt.Sprintf(
		`SELECT %s, %s FROM %s INNER JOIN %s  USING (%s) WHERE %s = ?`,
		playerIDColName, balanceColNmae, playersTableName, p2tTableName,
		playerIDColName, tournamentIDColName))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedSelectBackers, err = conn.Prepare(fmt.Sprintf(
		`SELECT %s.%s, %s FROM %s INNER JOIN %s ON %s.%s  = %s.%s WHERE %s = ? AND %s.%s = ?`,
		playersTableName, playerIDColName, balanceColNmae, playersTableName, p2bTableName, p2bTableName,
		backerColName, playersTableName, playerIDColName, tournamentIDColName, p2bTableName, playerIDColName))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	return tournamentConn, nil
}

func (tC *TournamentConnector) Close() {
	if tC.insert != nil {
		tC.insert.Close()
	}
	if tC.preparedUpdateTournament != nil {
		tC.preparedUpdateTournament.Close()
	}
	if tC.preparedSelectTournament != nil {
		tC.preparedSelectTournament.Close()
	}
	if tC.preparedJoinTournament != nil {
		tC.preparedJoinTournament.Close()
	}

	if tC.preparedInsertBackPlayer != nil {
		tC.preparedInsertBackPlayer.Close()
	}
	if tC.preparedSelectParticipants != nil {
		tC.preparedSelectParticipants.Close()
	}
	if tC.preparedSelectBackers != nil {
		tC.preparedSelectBackers.Close()
	}
}

func (tC *TournamentConnector) Create(tournament *entity.Tournament) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	_, err := tC.insert.Exec(tournament.ID(), tournament.Deposit())
	if err == nil {
		return nil
	}
	if err := err.(sqlite3.Error); err.Code == sqlite3.ErrConstraint {
		return ErrTournamentAlreadyExists
	}
	return err
}

func (tC *TournamentConnector) Read(id string) *entity.Tournament {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	var tournamentID string
	var tournamentDeposit float32
	var tournamentIsFinished bool
	err := tC.preparedSelectTournament.QueryRow(id).Scan(&tournamentID, &tournamentDeposit, &tournamentIsFinished)
	if err != nil {
		return nil
	}
	return entity.NewTournament(tournamentID, tournamentDeposit, tournamentIsFinished)
}

func (tC *TournamentConnector) Update(tournaments ...*entity.Tournament) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	preparedUpdateTx := tx.Stmt(tC.preparedUpdateTournament)
	for _, t := range tournaments {
		_, err := preparedUpdateTx.Exec(t.IsFinished(), t.ID())
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (tC *TournamentConnector) UpdateWithTx(tx *sql.Tx, tournaments ...*entity.Tournament) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	preparedUpdateTx := tx.Stmt(tC.preparedUpdateTournament)
	for _, t := range tournaments {
		_, err := preparedUpdateTx.Exec(t.IsFinished(), t.ID())
		if err != nil {
			return err
		}
	}
	return nil
}

func (tC *TournamentConnector) UpdateResult(tournament *entity.Tournament, players []*entity.Player) error {
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
	err = tC.UpdateWithTx(tx, tournament)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (tC *TournamentConnector) SelectBackPlayers(tournamentID, playerID string) ([]*entity.Player, error) {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	var backers []*entity.Player
	rows, err := tC.preparedSelectBackers.Query(tournamentID, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var backerID string
	var balance float32
	for rows.Next() {
		rows.Scan(&backerID, &balance)
		backers = append(backers, entity.NewPlayer(backerID, balance))
	}
	return backers, rows.Err()
}

func (tC *TournamentConnector) SelectParticipants(tournamentID string) ([]*entity.Player, [][]*entity.Player, error) {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	var participants []*entity.Player
	var backPlayers [][]*entity.Player

	rows, err := tC.preparedSelectParticipants.Query(tournamentID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var playerID string
	var balance float32
	for rows.Next() {
		rows.Scan(&playerID, &balance)
		backers, err := tC.SelectBackPlayers(tournamentID, playerID)
		if err != nil {
			return nil, nil, err
		}
		participants = append(participants, entity.NewPlayer(playerID, balance))
		backPlayers = append(backPlayers, backers)
	}
	return participants, backPlayers, rows.Err()
}

func (tC *TournamentConnector) UpdateJoin(tournament *entity.Tournament, player *entity.Player, backers []*entity.Player) error {
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

	_, err = tx.Stmt(tC.preparedJoinTournament).Exec(player.ID(), tournament.ID())
	if err != nil {
		tx.Rollback()
		return err
	}
	txPreparedBackPlayer := tx.Stmt(tC.preparedInsertBackPlayer)
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
