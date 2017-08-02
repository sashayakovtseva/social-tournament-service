package controller

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/sashayakovtseva/social-tournament-service/model"
)

var (
	playerControllerSingleton     *PlayerController
	playerControllerError         error
	playerControllerOnce sync.Once
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
	playerController.preparedInsert, err = connector.db.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME, BALANCE_COL_NAME))
	if err != nil {
		return nil, err
	}
	playerController.preparedUpdate, err = connector.db.Prepare(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?`,
		PLAYERS_TABLE_NAME, BALANCE_COL_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerController.preparedInsert.Close()
		return nil, err
	}
	playerController.preparedSelect, err = connector.db.Prepare(fmt.Sprintf(`SELECT %s, %s FROM %s WHERE %s = ?`,
		PLAYER_ID_COL_NAME, BALANCE_COL_NAME, PLAYERS_TABLE_NAME, PLAYER_ID_COL_NAME))
	if err != nil {
		playerController.preparedInsert.Close()
		playerController.preparedUpdate.Close()
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


func (pC *PlayerController) GetPlayerById(id string) *model.Player {
	player := &model.Player{}
	err := pC.preparedSelect.QueryRow(id).Scan(&player.Id, &player.Balance)
	if err != nil {
		return nil
	}
	return player
}

func (pC *PlayerController) UpdatePlayer(player *model.Player) error {
	_, err := pC.preparedUpdate.Exec(player.Balance, player.Id)
	return err
}


func (pC *PlayerController) InsertPlayer(player *model.Player) error {
	_, err := pC.preparedInsert.Exec(player.Id, player.Balance)
	return err
}