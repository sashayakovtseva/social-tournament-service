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
	PlayerConn             *PlayerConnector
	ErrNoSuchPlayer        = errors.New("no such player")
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrNotEnoughPoints     = errors.New("not enough points")
)

type PlayerConnector struct {
	insert *sql.Stmt
	update *sql.Stmt
	slct   *sql.Stmt
	take   *sql.Stmt
	fund   *sql.Stmt
}

func newPlayerConnector() (*PlayerConnector, error) {
	var err error
	playerConnector := new(PlayerConnector)
	playerConnector.insert, err = conn.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME, BALANCE_COL_NAME))
	if err != nil {
		playerConnector.Close()
		return nil, err
	}
	playerConnector.update, err = conn.Prepare(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?`,
		PLAYERS_TABLE_NAME, BALANCE_COL_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerConnector.Close()
		return nil, err
	}
	playerConnector.fund, err = conn.Prepare(fmt.Sprintf(`UPDATE %s SET %s = %s + ? WHERE %s = ?`,
		PLAYERS_TABLE_NAME, BALANCE_COL_NAME, BALANCE_COL_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerConnector.Close()
		return nil, err
	}
	playerConnector.take, err = conn.Prepare(fmt.Sprintf(`UPDATE %s SET %s = %s - ? WHERE %s = ?`,
		PLAYERS_TABLE_NAME, BALANCE_COL_NAME, BALANCE_COL_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerConnector.Close()
		return nil, err
	}
	playerConnector.slct, err = conn.Prepare(fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s = ?`,
		PLAYER_ID_COL_NAME, BALANCE_COL_NAME, PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerConnector.Close()
		return nil, err
	}

	return playerConnector, nil
}

func (pC *PlayerConnector) Close() {
	if pC.insert != nil {
		pC.insert.Close()
	}
	if pC.slct != nil {
		pC.slct.Close()
	}
	if pC.update != nil {
		pC.update.Close()
	}
	if pC.take != nil {
		pC.take.Close()
	}
	if pC.fund != nil {
		pC.fund.Close()
	}
}

func (pC *PlayerConnector) Create(player *entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	_, err := pC.insert.Exec(player.Id(), player.Balance())
	if err == nil {
		return nil
	}
	if err := err.(sqlite3.Error); err.Code == sqlite3.ErrConstraint {
		return ErrPlayerAlreadyExists
	}
	return err
}

func (pC *PlayerConnector) Read(id string) *entity.Player {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	var playerId string
	var balance float32
	err := pC.slct.QueryRow(id).Scan(&playerId, &balance)
	if err != nil {
		return nil
	}
	return entity.NewPlayer(playerId, balance)
}

func (pC *PlayerConnector) Take(playerId string, points float32) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	r, err := pC.take.Exec(points, playerId)
	if err != nil {
		if err := err.(sqlite3.Error); err.Code == sqlite3.ErrConstraint {
			return ErrNotEnoughPoints
		}
		return err
	}

	if n, _ := r.RowsAffected(); n == 0 {
		return ErrNoSuchPlayer
	}
	return nil
}

func (pC *PlayerConnector) Fund(playerId string, points float32) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	r, err := pC.fund.Exec(points, playerId)
	if err != nil {
		return err
	}
	if n, _ := r.RowsAffected(); n == 0 {
		return ErrNoSuchPlayer
	}
	return nil
}

func (pC *PlayerConnector) UpdateWithTx(tx *sql.Tx, players ...*entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	preparedUpdateTx := tx.Stmt(pC.update)
	for _, player := range players {
		_, err := preparedUpdateTx.Exec(player.Balance(), player.Id())
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	var err error
	PlayerConn, err = newPlayerConnector()
	if err != nil {
		log.Fatal(err)
	}
}
