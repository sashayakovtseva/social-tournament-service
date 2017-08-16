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
		TOURNAMENTS_TABLE_NAME, TOURNAMENT_ID_COL_NAME, DEPOSIT_COL_NAME))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedUpdateTournament, err = conn.Prepare(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?`,
		TOURNAMENTS_TABLE_NAME, FINISHED_COL_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedSelectTournament, err = conn.Prepare(fmt.Sprintf(
		`SELECT %s, %s, %s FROM %s WHERE %s = ?`,
		TOURNAMENT_ID_COL_NAME, DEPOSIT_COL_NAME, FINISHED_COL_NAME, TOURNAMENTS_TABLE_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedJoinTournament, err = conn.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		P2T_TABLE_NAME, PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedInsertBackPlayer, err = conn.Prepare(fmt.Sprintf(
		`INSERT INTO %s(%s,%s, %s) values(?,?, ?)`,
		P2B_TABLE_NAME, PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME, BACKER_COL_NAME))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedSelectParticipants, err = conn.Prepare(fmt.Sprintf(
		`SELECT %s, %s FROM %s INNER JOIN %s  USING (%s) WHERE %s = ?`,
		PLAYER_ID_COL_NAME, BALANCE_COL_NAME, PLAYERS_TABLE_NAME, P2T_TABLE_NAME,
		PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournamentConn.Close()
		return nil, err
	}
	tournamentConn.preparedSelectBackers, err = conn.Prepare(fmt.Sprintf(
		`SELECT %s.%s, %s FROM %s INNER JOIN %s ON %s.%s  = %s.%s WHERE %s = ? AND %s.%s = ?`,
		PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME, BALANCE_COL_NAME, PLAYERS_TABLE_NAME, P2B_TABLE_NAME, P2B_TABLE_NAME,
		BACKER_COL_NAME, PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME, P2B_TABLE_NAME, PLAYER_ID_COL_NAME))
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

	_, err := tC.insert.Exec(tournament.Id(), tournament.Deposit())
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

	var tournamentId string
	var tournamentDeposit float32
	var tournamentIsFinished bool
	err := tC.preparedSelectTournament.QueryRow(id).Scan(&tournamentId, &tournamentDeposit, &tournamentIsFinished)
	if err != nil {
		return nil
	}
	return entity.NewTournament(tournamentId, tournamentDeposit, tournamentIsFinished)
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
		_, err := preparedUpdateTx.Exec(t.IsFinished(), t.Id())
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
		_, err := preparedUpdateTx.Exec(t.IsFinished(), t.Id())
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

func (tC *TournamentConnector) SelectBackPlayers(tournamentId, playerId string) ([]*entity.Player, error) {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	var backers []*entity.Player
	rows, err := tC.preparedSelectBackers.Query(tournamentId, playerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var backerId string
	var balance float32
	for rows.Next() {
		rows.Scan(&backerId, &balance)
		backers = append(backers, entity.NewPlayer(backerId, balance))
	}
	return backers, rows.Err()
}

func (tC *TournamentConnector) SelectParticipants(tournamentId string) ([]*entity.Player, [][]*entity.Player, error) {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	var participants []*entity.Player
	var backPlayers [][]*entity.Player

	rows, err := tC.preparedSelectParticipants.Query(tournamentId)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var playerId string
	var balance float32
	for rows.Next() {
		rows.Scan(&playerId, &balance)
		backers, err := tC.SelectBackPlayers(tournamentId, playerId)
		if err != nil {
			return nil, nil, err
		}
		participants = append(participants, entity.NewPlayer(playerId, balance))
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

	_, err = tx.Stmt(tC.preparedJoinTournament).Exec(player.Id(), tournament.Id())
	if err != nil {
		tx.Rollback()
		return err
	}
	txPreparedBackPlayer := tx.Stmt(tC.preparedInsertBackPlayer)
	for _, backer := range backers {
		_, err := txPreparedBackPlayer.Exec(player.Id(), tournament.Id(), backer.Id())
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
