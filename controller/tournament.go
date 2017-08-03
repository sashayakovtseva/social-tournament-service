package controller

import (
	"database/sql"
	"fmt"
	"sync"
	"errors"

	"github.com/sashayakovtseva/social-tournament-service/entity"
)

var (
	tournControllerSingleton *TournamentController
	tournControllerError     error
	tournControllerOnce      sync.Once
)

type TournamentController struct {
	connector              *DBConnector
	preparedInsert         *sql.Stmt
	preparedUpdate         *sql.Stmt
	preparedSelect         *sql.Stmt
	preparedJoinTournament *sql.Stmt
	preparedBackPlayer     *sql.Stmt
}

func newTournamentController() (*TournamentController, error) {
	connector, err := GetConnector()
	if err != nil {
		return nil, err
	}

	tournController := &TournamentController{connector: connector}
	tournController.preparedInsert, err = connector.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		TOURNAMENTS_TABLE_NAME, TOURNAMENT_ID_COL_NAME, DEPOSIT_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedUpdate, err = connector.Prepare(fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s = ?`,
		TOURNAMENTS_TABLE_NAME, WINNER_COL_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedSelect, err = connector.Prepare(fmt.Sprintf(`SELECT %s, %s, %s FROM %s WHERE %s = ?`,
		TOURNAMENT_ID_COL_NAME, DEPOSIT_COL_NAME, WINNER_COL_NAME, TOURNAMENTS_TABLE_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedJoinTournament, err = connector.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s) values(?,?)`,
		P2T_TABLE_NAME, PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	tournController.preparedBackPlayer, err = connector.Prepare(fmt.Sprintf(`INSERT INTO %s(%s,%s, %s) values(?,?, ?)`,
		P2B_TABLE_NAME, PLAYER_ID_COL_NAME, TOURNAMENT_ID_COL_NAME, BACKER_COL_NAME))
	if err != nil {
		tournController.Close()
		return nil, err
	}
	return tournController, nil
}

func GetTournamentController() (*TournamentController, error) {
	tournControllerOnce.Do(func() {
		tournControllerSingleton, tournControllerError = newTournamentController()
	})
	return tournControllerSingleton, tournControllerError
}

func (tC *TournamentController) Close() {
	if tC.preparedUpdate != nil {
		tC.preparedUpdate.Close()
	}
	if tC.preparedInsert != nil {
		tC.preparedInsert.Close()
	}
	if tC.preparedSelect != nil {
		tC.preparedSelect.Close()
	}
	if tC.preparedJoinTournament != nil {
		tC.preparedJoinTournament.Close()
	}
}

func (tC *TournamentController) GetById(id string) *entity.Tournament {
	tournament := &entity.Tournament{}
	err := tC.preparedSelect.QueryRow(id).Scan(&tournament.Id, &tournament.Deposit, &tournament.WinnerId)
	if err != nil {
		return nil
	}
	return tournament
}

func (tC *TournamentController) Update(tournament *entity.Tournament) error {
	_, err := tC.preparedUpdate.Exec(tournament.WinnerId, tournament.Id)
	return err
}

func (tC *TournamentController) Create(tournament *entity.Tournament) error {
	_, err := tC.preparedInsert.Exec(tournament.Id, tournament.Deposit)
	return err
}

func (tC *TournamentController) Announce() error {
 return  nil
}

func (tC *TournamentController) JoinTournament(tournament *entity.Tournament, player *entity.Player, backers []*entity.Player) error {
	if !tournament.Join(player, backers) {
		return errors.New("Not enough points to join")
	}

	tx, err := tC.connector.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Stmt(tC.preparedJoinTournament).Exec(player.Id, tournament.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	txPreparedBackPlayer := tx.Stmt(tC.preparedBackPlayer)
	for _, backer := range backers {
		_, err := txPreparedBackPlayer.Exec(player.Id, tournament.Id, backer.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
