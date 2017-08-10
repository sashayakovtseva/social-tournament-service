package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	PlayerConn *PlayerConnector
)

type PlayerConnector struct {
	preparedInsert *sql.Stmt
	preparedUpdate *sql.Stmt
	preparedSelect *sql.Stmt
}

func newPlayerConnector() (*PlayerConnector, error) {
	var err error
	playerConnector := new(PlayerConnector)
	playerConnector.preparedInsert, err = conn.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME, BALANCE_COL_NAME))
	if err != nil {
		playerConnector.Close()
		return nil, err
	}
	playerConnector.preparedUpdate, err = conn.Prepare(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?`,
		PLAYERS_TABLE_NAME, BALANCE_COL_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerConnector.Close()
		return nil, err
	}
	playerConnector.preparedSelect, err = conn.Prepare(fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s = ?`,
		PLAYER_ID_COL_NAME, BALANCE_COL_NAME, PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerConnector.Close()
		return nil, err
	}
	return playerConnector, nil
}

func (pC *PlayerConnector) Close() {
	if pC.preparedInsert != nil {
		pC.preparedInsert.Close()
	}
	if pC.preparedSelect != nil {
		pC.preparedSelect.Close()
	}
	if pC.preparedUpdate != nil {
		pC.preparedUpdate.Close()
	}
}

func (pC *PlayerConnector) Create(player *entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	_, err := pC.preparedInsert.Exec(player.Id(), player.Balance())
	return err
}

func (pC *PlayerConnector) Read(id string) *entity.Player {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	var playerId string
	var balance float32
	err := pC.preparedSelect.QueryRow(id).Scan(&playerId, &balance)
	if err != nil {
		return nil
	}
	return entity.NewPlayer(playerId, balance)
}

func (pC *PlayerConnector) Update(players ...*entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	preparedUpdateTx := tx.Stmt(pC.preparedUpdate)
	for _, player := range players {
		_, err := preparedUpdateTx.Exec(player.Balance(), player.Id())
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (pC *PlayerConnector) UpdateWithTx(tx *sql.Tx, players ...*entity.Player) error {
	resetMutex.RLock()
	defer resetMutex.RUnlock()

	preparedUpdateTx := tx.Stmt(pC.preparedUpdate)
	for _, player := range players {
		_, err := preparedUpdateTx.Exec(player.Balance(), player.Id())
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	log.Println("init player.go")
	var err error
	PlayerConn, err = newPlayerConnector()
	if err != nil {
		log.Fatal(err)
	}
}
