package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	playerControllerSingleton *PlayerController
	playerControllerError     error
	playerControllerOnce      sync.Once
)

type PlayerController struct {
	connector      *DBConnector
	preparedInsert *sql.Stmt
	preparedUpdate *sql.Stmt
	preparedSelect *sql.Stmt
}

func newPlayerController() (*PlayerController, error) {
	connector, err := GetConnector()
	if err != nil {
		return nil, err
	}

	playerController := &PlayerController{connector: connector}
	playerController.preparedInsert, err = connector.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME, BALANCE_COL_NAME))
	if err != nil {
		playerController.Close()
		return nil, err
	}
	playerController.preparedUpdate, err = connector.Prepare(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?`,
		PLAYERS_TABLE_NAME, BALANCE_COL_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerController.Close()
		return nil, err
	}
	playerController.preparedSelect, err = connector.Prepare(fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s = ?`,
		PLAYER_ID_COL_NAME, BALANCE_COL_NAME, PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerController.Close()
		return nil, err
	}
	return playerController, nil
}

func GetPlayerController() (*PlayerController, error) {
	playerControllerOnce.Do(func() {
		playerControllerSingleton, playerControllerError = newPlayerController()
	})
	return playerControllerSingleton, playerControllerError
}

func (pC *PlayerController) Close() {
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

func (pC *PlayerController) GetById(id string) *entity.Player {
	player := &entity.Player{}
	err := pC.preparedSelect.QueryRow(id).Scan(&player.Id, &player.Balance)
	if err != nil {
		return nil
	}
	return player
}

func (pC *PlayerController) Update(player *entity.Player) error {
	_, err := pC.preparedUpdate.Exec(player.Balance, player.Id)
	return err
}

func (pC *PlayerController) Create(player *entity.Player) error {
	_, err := pC.preparedInsert.Exec(player.Id, player.Balance)
	return err
}

func (pC *PlayerController) Take(playerId string, points int) error {
	player := pC.GetById(playerId)
	if player == nil {
		return errors.New("No such player")
	}
	if player.Take(points) {
		pC.Update(player)
		return nil
	}
	return errors.New("Not enough points")
}

func (pC *PlayerController) Fund(playerId string, points int) error {
	player := pC.GetById(playerId)
	if player == nil {
		return pC.Create(entity.NewPlayer(playerId, points))
	}
	player.Fund(points)
	return pC.Update(player)
}
