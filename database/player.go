package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	PlayerConn             *PlayerConnector
	ErrNoSuchPlayer        = errors.New("no such player")
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrNotEnoughPoints     = errors.New("not enough points")
)

type PlayerConnector struct {
	insert     *sql.Stmt
	update     *sql.Stmt
	read       *sql.Stmt
	take       *sql.Stmt
	fund       *sql.Stmt
	statements []*sql.Stmt
}

func newPlayerConnector() (*PlayerConnector, error) {
	var err error

	playerConnector := new(PlayerConnector)

	playerConnector.insert, err = prepareAndAdd(playerConnector.statements, `INSERT INTO %s(%s,%s) values(?,?)`,
		playersTableName, playerIDColName, balanceColName)
	if haveToFailAndClose(playerConnector, err) {
		return nil, err
	}

	playerConnector.update, err = prepareAndAdd(playerConnector.statements, `UPDATE %s SET %s = ? WHERE %s = ?`,
		playersTableName, balanceColName, playerIDColName)
	if haveToFailAndClose(playerConnector, err) {
		return nil, err
	}

	playerConnector.fund, err = prepareAndAdd(playerConnector.statements, `UPDATE %s SET %s = %s + ? WHERE %s = ?`,
		playersTableName, balanceColName, balanceColName, playerIDColName)
	if haveToFailAndClose(playerConnector, err) {
		return nil, err
	}

	playerConnector.take, err = prepareAndAdd(playerConnector.statements, `UPDATE %s SET %s = %s - ? WHERE %s = ?`,
		playersTableName, balanceColName, balanceColName, playerIDColName)
	if haveToFailAndClose(playerConnector, err) {
		return nil, err
	}

	playerConnector.read, err = prepareAndAdd(playerConnector.statements, `SELECT %s, %s FROM %s WHERE %s = ?`,
		playerIDColName, balanceColName, playersTableName, playerIDColName)
	if haveToFailAndClose(playerConnector, err) {
		return nil, err
	}

	return playerConnector, nil
}

func (connector *PlayerConnector) Close() {
	for _, stmt := range connector.statements {
		checkAndClose(stmt)
	}
}

func (connector *PlayerConnector) Create(player *entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	_, err := connector.insert.Exec(player.ID(), player.Balance())
	return replaceConstraintWithCustom(err, ErrPlayerAlreadyExists)
}

func (connector *PlayerConnector) Read(id string) *entity.Player {
	var playerID string
	var balance float32

	resetMutex.RLock()
	defer resetMutex.RUnlock()

	err := connector.read.QueryRow(id).Scan(&playerID, &balance)
	if err != nil {
		return nil
	}
	return entity.NewPlayer(playerID, balance)
}

func (connector *PlayerConnector) TakePoints(playerID string, points float32) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	res, err := connector.take.Exec(points, playerID)
	if err != nil {
		return replaceConstraintWithCustom(err, ErrNotEnoughPoints)
	}
	return checkPlayerExists(res)
}

func (connector *PlayerConnector) FundPoints(playerID string, points float32) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	res, err := connector.fund.Exec(points, playerID)
	if err != nil {
		return err
	}
	return checkPlayerExists(res)
}

func checkPlayerExists(result sql.Result) error {
	if n, _ := result.RowsAffected(); n == 0 {
		return ErrNoSuchPlayer
	}
	return nil
}

func (connector *PlayerConnector) UpdateWithTx(tx *sql.Tx, players ...*entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	preparedUpdateTx := tx.Stmt(connector.update)
	for _, player := range players {
		_, err := preparedUpdateTx.Exec(player.Balance(), player.ID())
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
